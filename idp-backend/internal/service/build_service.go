package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"idp-backend/internal/domain"
)

type BuildService struct {
	Registry   string
	k8sService *K8sService
}

func NewBuildService(registry string, k8sService *K8sService) *BuildService {
	return &BuildService{
		Registry:   registry,
		k8sService: k8sService,
	}
}

func (s *BuildService) Build(project *domain.Project) (*domain.BuildResult, error) {
	// 1. Create temp directory
	tempDir, err := os.MkdirTemp("", "idp-build-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir)

	// 2. Clone repository
	log.Printf("[BUILD] Cloning %s to %s", project.GitURL, tempDir)
	cmd := exec.Command("git", "clone", project.GitURL, ".")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to clone repo: %v", err)
	}

	// 3. Detect framework
	framework, err := s.DetectFramework(tempDir)
	if err != nil {
		return nil, err
	}
	log.Printf("[BUILD] Detected framework: %s", framework)

	// 4. Generate Dockerfile
	dockerfile := s.GenerateDockerfile(framework)
	err = os.WriteFile(filepath.Join(tempDir, "Dockerfile"), []byte(dockerfile), 0644)
	if err != nil {
		return nil, err
	}

	// 5. Build and Push Image (using Kaniko Job in K8s)
	imageName := fmt.Sprintf("%s/%s:latest", s.Registry, project.Name)
	log.Printf("[BUILD] Starting Kaniko build job for: %s", imageName)
	
	ctx := context.Background()
	if err := s.k8sService.CreateBuildJob(ctx, project, dockerfile, imageName); err != nil {
		return nil, fmt.Errorf("kaniko build failed: %v", err)
	}

	log.Printf("[BUILD] Kaniko build successful: %s", imageName)

	// 6. Return Result
	port := 80
	if framework == domain.FrameworkNextSSR {
		port = 3000
	}

	return &domain.BuildResult{
		Framework: framework,
		ImageName: imageName,
		Port:      port,
	}, nil
}

func (s *BuildService) DetectFramework(path string) (domain.FrameworkType, error) {
	// 1. Next.js check (explicit config files)
	if _, err := os.Stat(filepath.Join(path, "next.config.js")); err == nil {
		return domain.FrameworkNextSSR, nil
	}
	if _, err := os.Stat(filepath.Join(path, "next.config.mjs")); err == nil {
		return domain.FrameworkNextSSR, nil
	}

	// 2. package.json check
	pkgPath := filepath.Join(path, "package.json")
	if _, err := os.Stat(pkgPath); err == nil {
		data, err := os.ReadFile(pkgPath)
		if err != nil {
			return domain.FrameworkUnknown, err
		}

		var pkg struct {
			Scripts      map[string]string `json:"scripts"`
			Dependencies map[string]string `json:"dependencies"`
		}
		if err := json.Unmarshal(data, &pkg); err != nil {
			return domain.FrameworkUnknown, err
		}

		// Check for static export
		if build, ok := pkg.Scripts["build"]; ok && strings.Contains(build, "export") {
			return domain.FrameworkStatic, nil
		}
		
		// Check for Next.js in dependencies
		if deps := pkg.Dependencies; deps != nil {
			if _, ok := deps["next"]; ok {
				return domain.FrameworkNextSSR, nil
			}
		}

		// Check for start script (likely SSR/Node app)
		if _, ok := pkg.Scripts["start"]; ok {
			return domain.FrameworkNextSSR, nil
		}
	}

	return domain.FrameworkStatic, nil // Fallback to static
}

func (s *BuildService) GenerateDockerfile(framework domain.FrameworkType) string {
	switch framework {
	case domain.FrameworkNextSSR:
		return `FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install --legacy-peer-deps
COPY . .
RUN npm run build

FROM node:20-alpine
WORKDIR /app
COPY --from=builder /app/package*.json ./
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/node_modules ./node_modules
EXPOSE 3000
CMD ["npm", "start"]`

	case domain.FrameworkStatic:
		return `FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install --legacy-peer-deps
COPY . .
RUN npm run build

# Validation and normalization Step
RUN if [ -d "dist" ] && [ -f "dist/index.html" ]; then \
    mv dist dist_output; \
elif [ -d "build" ] && [ -f "build/index.html" ]; then \
    mv build dist_output; \
elif [ -d "out" ] && [ -f "out/index.html" ]; then \
    mv out dist_output; \
else \
    echo "Build succeeded but no valid output found (dist/build/out index.html missing)"; \
    exit 1; \
fi

FROM nginx:alpine
COPY --from=builder /app/dist_output /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]`

	default:
		return `FROM nginx:alpine
COPY . /usr/share/nginx/html
EXPOSE 80`
	}
}
