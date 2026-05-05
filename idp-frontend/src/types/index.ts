export type ProjectStatus = "PENDING" | "BUILDING" | "READY" | "FAILED";

export interface Project {
  id: string;
  name: string;
  git_url: string;
  domain: string;
  status: ProjectStatus;
  branch: string;
  created_at: string;
  updated_at: string;
}

export interface DeploymentLog {
  id: string;
  project_id: string;
  content: string;
  timestamp: string;
}
