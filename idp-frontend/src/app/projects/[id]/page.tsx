"use client";

import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import { 
  ArrowLeft, 
  ExternalLink, 
  GitPullRequest, 
  GitBranch, 
  Clock, 
  Terminal, 
  Globe, 
  Settings, 
  Trash2,
  AlertTriangle
} from "lucide-react";
import Link from "next/link";
import { Project, DeploymentLog } from "@/types";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { formatDistanceToNow } from "date-fns";

export default function ProjectDetail() {
  const { id } = useParams();
  const router = useRouter();
  const [project, setProject] = useState<Project | null>(null);
  const [logs, setLogs] = useState<DeploymentLog[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const apiUrl = process.env.NEXT_PUBLIC_API_URL || "";
        const [projRes, logsRes] = await Promise.all([
          fetch(`${apiUrl}/api/projects/${id}`),
          fetch(`${apiUrl}/api/projects/${id}/logs`)
        ]);

        if (projRes.ok) {
          const projData = await projRes.json();
          setProject(projData);
        }
        
        if (logsRes.ok) {
          const logsData = await logsRes.json();
          setLogs(Array.isArray(logsData) ? logsData : []);
        }
      } catch (error) {
        console.error("Failed to fetch project data:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
    const interval = setInterval(fetchData, 5000);
    return () => clearInterval(interval);
  }, [id]);

  if (loading) {
    return (
      <div className="mx-auto max-w-7xl px-4 py-20 animate-pulse">
        <div className="h-8 w-48 bg-gray-200 rounded mb-8" />
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          <div className="lg:col-span-2 space-y-8">
            <div className="h-96 bg-gray-200 rounded-2xl" />
            <div className="h-64 bg-gray-200 rounded-2xl" />
          </div>
          <div className="space-y-8">
            <div className="h-48 bg-gray-200 rounded-2xl" />
          </div>
        </div>
      </div>
    );
  }

  if (!project) {
    return (
      <div className="flex flex-col items-center justify-center py-32">
        <AlertTriangle className="h-12 w-12 text-amber-500 mb-4" />
        <h1 className="text-2xl font-bold">Project not found</h1>
        <Link href="/dashboard" className="mt-4 text-blue-500 hover:underline">
          Return to dashboard
        </Link>
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
      <div className="mb-10 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-4">
          <Link
            href="/dashboard"
            className="flex h-10 w-10 items-center justify-center rounded-full border border-gray-200 bg-white transition-colors hover:bg-gray-50 dark:border-gray-800 dark:bg-black dark:hover:bg-gray-900"
          >
            <ArrowLeft className="h-5 w-5" />
          </Link>
          <div>
            <h1 className="text-3xl font-bold tracking-tight text-gray-900 dark:text-white">
              {project.name}
            </h1>
            <div className="mt-1 flex items-center gap-3">
              <StatusBadge status={project.status} />
              <span className="text-sm text-gray-500">
                Last updated {project.updated_at && !isNaN(new Date(project.updated_at).getTime())
                  ? `${formatDistanceToNow(new Date(project.updated_at))} ago`
                  : "recently"}
              </span>
            </div>
          </div>
        </div>

        <div className="flex items-center gap-3">
          <a
            href={project.domain ? `https://${project.domain}` : "#"}
            target="_blank"
            rel="noopener noreferrer"
            className="flex items-center gap-2 rounded-lg bg-black px-6 py-2.5 text-sm font-medium text-white transition-opacity hover:opacity-90 dark:bg-white dark:text-black"
          >
            <ExternalLink className="h-4 w-4" />
            Visit Site
          </a>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-8 lg:grid-cols-3">
        {/* Main Content */}
        <div className="lg:col-span-2 space-y-8">
          {/* Site Preview */}
          <div className="overflow-hidden rounded-2xl border border-gray-200 bg-white shadow-sm dark:border-gray-800 dark:bg-black">
            <div className="border-b border-gray-100 bg-gray-50 px-6 py-3 dark:border-gray-800 dark:bg-gray-900/50">
              <div className="flex items-center gap-2">
                <Globe className="h-4 w-4 text-gray-400" />
                <span className="text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Deployment Preview
                </span>
              </div>
            </div>
            <div className="relative aspect-video w-full bg-gray-100 dark:bg-gray-900">
              {project.status === "READY" ? (
                <div className="flex h-full w-full items-center justify-center">
                   <div className="text-center">
                     <Globe className="h-12 w-12 text-gray-300 mx-auto mb-4" />
                     <p className="text-sm text-gray-400 font-medium">Preview available at</p>
                     <p className="text-lg text-gray-900 dark:text-white font-bold">{project.domain}</p>
                   </div>
                </div>
              ) : (
                <div className="flex h-full w-full flex-col items-center justify-center">
                  <div className="h-12 w-12 animate-pulse rounded-full bg-gray-200 dark:bg-gray-800" />
                  <p className="mt-4 text-sm font-medium text-gray-500">
                    Deployment in progress...
                  </p>
                </div>
              )}
            </div>
          </div>

          {/* Logs */}
          <div className="rounded-2xl border border-gray-900 bg-[#0A0A0A] overflow-hidden shadow-2xl">
            <div className="flex items-center justify-between border-b border-white/10 bg-white/5 px-6 py-4">
              <div className="flex items-center gap-2">
                <Terminal className="h-4 w-4 text-emerald-400" />
                <span className="text-sm font-bold text-white uppercase tracking-tight">
                  Deployment Logs
                </span>
              </div>
              <div className="flex gap-1.5">
                <div className="h-2 w-2 rounded-full bg-rose-500/50" />
                <div className="h-2 w-2 rounded-full bg-amber-500/50" />
                <div className="h-2 w-2 rounded-full bg-emerald-500/50" />
              </div>
            </div>
            <div className="h-80 overflow-y-auto p-6 font-mono text-[13px] leading-relaxed text-gray-300">
              {logs.length > 0 ? (
                logs.map((log, index) => (
                  <div key={log.id} className="mb-1 flex gap-4">
                    <span className="shrink-0 select-none text-gray-600 w-24">
                      {new Date(log.timestamp).toLocaleTimeString([], { hour12: false })}
                    </span>
                    <span className="whitespace-pre-wrap">{log.content}</span>
                  </div>
                ))
              ) : (
                <div className="flex h-full items-center justify-center text-gray-600 italic">
                  No logs available for this deployment yet.
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Sidebar */}
        <div className="space-y-8">
          <div className="rounded-2xl border border-gray-200 bg-white p-6 dark:border-gray-800 dark:bg-black">
            <h2 className="mb-6 text-sm font-bold uppercase tracking-tight text-gray-900 dark:text-white">
              Deployment Info
            </h2>
            <div className="space-y-6">
              <div>
                <label className="text-[11px] font-bold uppercase tracking-wider text-gray-500">
                  Domain
                </label>
                <div className="mt-1 flex items-center gap-2 font-medium text-gray-900 dark:text-white">
                  <Globe className="h-4 w-4 text-gray-400" />
                  {project.domain || "Not assigned"}
                </div>
              </div>
              <div>
                <label className="text-[11px] font-bold uppercase tracking-wider text-gray-500">
                  Repository
                </label>
                <div className="mt-1 flex items-center gap-2 font-medium text-gray-900 dark:text-white">
                  <GitPullRequest className="h-4 w-4 text-gray-400" />
                  <span className="truncate">{project.git_url.split('/').slice(-2).join('/')}</span>
                </div>
              </div>
              <div>
                <label className="text-[11px] font-bold uppercase tracking-wider text-gray-500">
                  Branch
                </label>
                <div className="mt-1 flex items-center gap-2 font-medium text-gray-900 dark:text-white">
                  <GitBranch className="h-4 w-4 text-gray-400" />
                  {project.branch || "main"}
                </div>
              </div>
            </div>
          </div>

          <div className="rounded-2xl border border-rose-200 bg-rose-50/50 p-6 dark:border-rose-900/30 dark:bg-rose-900/10">
            <h2 className="mb-4 text-sm font-bold uppercase tracking-tight text-rose-900 dark:text-rose-400">
              Danger Zone
            </h2>
            <p className="mb-6 text-xs text-rose-700/80 dark:text-rose-400/60 leading-relaxed">
              Once you delete a project, there is no going back. Please be certain.
            </p>
            <button className="flex w-full items-center justify-center gap-2 rounded-lg bg-rose-600 px-4 py-2 text-sm font-bold text-white transition-opacity hover:opacity-90">
              <Trash2 className="h-4 w-4" />
              Delete Project
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
