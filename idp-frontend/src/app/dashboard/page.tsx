"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { Search, Plus, SlidersHorizontal, RefreshCcw } from "lucide-react";
import { Project } from "@/types";
import { ProjectCard } from "@/components/ui/ProjectCard";

export default function Dashboard() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");

  const fetchProjects = async () => {
    setLoading(true);
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/projects`,
      );
      if (response.ok) {
        const data = await response.json();
        setProjects(data);
      }
    } catch (error) {
      console.error("Failed to fetch projects:", error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchProjects();
  }, []);

  const filteredProjects = (projects || []).filter((p) =>
    p.name.toLowerCase().includes(search.toLowerCase()),
  );

  return (
    <div className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
      <div className="mb-10 flex flex-col items-start justify-between gap-6 sm:flex-row sm:items-center">
        <div className="relative w-full max-w-md">
          <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
            <Search className="h-4 w-4 text-gray-400" />
          </div>
          <input
            type="text"
            placeholder="Search projects..."
            className="block w-full rounded-lg border border-gray-200 bg-white py-2.5 pl-10 pr-3 text-sm transition-all focus:border-gray-900 focus:ring-1 focus:ring-gray-900 dark:border-gray-800 dark:bg-black dark:focus:border-white dark:focus:ring-white"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>

        <div className="flex w-full items-center gap-3 sm:w-auto">
          <button
            onClick={fetchProjects}
            className="flex items-center gap-2 rounded-lg border border-gray-200 bg-white px-4 py-2.5 text-sm font-medium transition-colors hover:bg-gray-50 dark:border-gray-800 dark:bg-black dark:hover:bg-gray-900"
          >
            <RefreshCcw
              className={`h-4 w-4 ${loading ? "animate-spin" : ""}`}
            />
            <span className="hidden sm:inline">Refresh</span>
          </button>
          <Link
            href="/new"
            className="flex flex-1 items-center justify-center gap-2 rounded-lg bg-black px-6 py-2.5 text-sm font-medium text-white transition-opacity hover:opacity-90 sm:flex-none dark:bg-white dark:text-black"
          >
            <Plus className="h-4 w-4" />
            New Project
          </Link>
        </div>
      </div>

      {loading ? (
        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {[...Array(6)].map((_, i) => (
            <div
              key={i}
              className="h-48 animate-pulse rounded-xl border border-gray-100 bg-white dark:border-gray-800 dark:bg-black"
            />
          ))}
        </div>
      ) : filteredProjects.length > 0 ? (
        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {filteredProjects.map((project) => (
            <ProjectCard key={project.id} project={project} />
          ))}
        </div>
      ) : (
        <div className="flex flex-col items-center justify-center rounded-2xl border border-dashed border-gray-200 py-24 dark:border-gray-800">
          <div className="rounded-full bg-gray-50 p-4 dark:bg-gray-900">
            <SlidersHorizontal className="h-8 w-8 text-gray-400" />
          </div>
          <h3 className="mt-4 text-lg font-medium text-gray-900 dark:text-white">
            No projects found
          </h3>
          <p className="mt-1 text-sm text-gray-500">
            {search
              ? "Try adjusting your search query."
              : "Get started by creating your first project."}
          </p>
          {!search && (
            <Link
              href="/new"
              className="mt-6 flex items-center gap-2 rounded-lg bg-black px-6 py-2.5 text-sm font-medium text-white transition-opacity hover:opacity-90 dark:bg-white dark:text-black"
            >
              <Plus className="h-4 w-4" />
              Deploy Project
            </Link>
          )}
        </div>
      )}
    </div>
  );
}
