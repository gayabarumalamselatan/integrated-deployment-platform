import Link from "next/link";
import { formatDistanceToNow } from "date-fns";
import { ExternalLink, GitBranch, Clock } from "lucide-react";
import { Project } from "@/types";
import { StatusBadge } from "./StatusBadge";
import { cn } from "@/lib/utils";

interface ProjectCardProps {
  project: Project;
}

export function ProjectCard({ project }: ProjectCardProps) {
  return (
    <Link
      href={`/projects/${project.id}`}
      className="group relative flex flex-col justify-between rounded-xl border border-gray-200 bg-white p-6 transition-all duration-300 hover:border-gray-900 hover:shadow-[0_8px_30px_rgb(0,0,0,0.12)] dark:border-gray-800 dark:bg-black dark:hover:border-white"
    >
      <div>
        <div className="flex items-start justify-between">
          <div className="space-y-1">
            <h3 className="text-lg font-bold tracking-tight text-gray-900 dark:text-white">
              {project.name}
            </h3>
            <div className="flex items-center gap-2 text-sm text-gray-500 dark:text-gray-400">
              <ExternalLink className="h-3 w-3" />
              <span className="hover:underline transition-colors">
                {project.domain || "no-domain-set.idp.dev"}
              </span>
            </div>
          </div>
          <StatusBadge status={project.status} />
        </div>

        <div className="mt-6 flex flex-wrap gap-4">
          <div className="flex items-center gap-1.5 text-xs text-gray-500 dark:text-gray-400">
            <GitBranch className="h-3.5 w-3.5" />
            <span className="font-medium text-gray-700 dark:text-gray-300">
              {project.branch || "main"}
            </span>
          </div>
          <div className="flex items-center gap-1.5 text-xs text-gray-500 dark:text-gray-400">
            <Clock className="h-3.5 w-3.5" />
            <span>
              {project.updated_at && !isNaN(new Date(project.updated_at).getTime())
                ? `${formatDistanceToNow(new Date(project.updated_at))} ago`
                : "Just now"}
            </span>
          </div>
        </div>
      </div>

      <div className="mt-6 flex items-center justify-between border-t border-gray-100 pt-4 dark:border-gray-800">
        <div className="flex items-center gap-2 text-xs font-medium text-gray-500 dark:text-gray-400">
          <GitBranch className="h-3.5 w-3.5" />
          <span className="truncate max-w-[150px]">{project.git_url}</span>
        </div>
        <div className="rounded-full bg-gray-50 p-1.5 opacity-0 transition-opacity group-hover:opacity-100 dark:bg-gray-900">
          <ExternalLink className="h-3.5 w-3.5 text-gray-900 dark:text-white" />
        </div>
      </div>
    </Link>
  );
}
