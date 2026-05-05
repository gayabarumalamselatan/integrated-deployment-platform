"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import {
  CodeXml,
  ArrowLeft,
  Loader2,
  Sparkles,
  AlertCircle,
} from "lucide-react";
import Link from "next/link";

export default function NewProject() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [formData, setFormData] = useState({
    name: "",
    git_url: "",
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/projects`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(formData),
      });

      if (response.ok) {
        router.push("/dashboard");
      } else {
        const data = await response.json();
        setError(data.message || "Failed to create project");
      }
    } catch (err) {
      setError(`Connection refused. Make sure your backend is running at ${process.env.NEXT_PUBLIC_API_URL}`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="mx-auto max-w-2xl px-4 py-20 sm:px-6 lg:px-8">
      <Link
        href="/dashboard"
        className="mb-8 inline-flex items-center gap-2 text-sm text-gray-500 hover:text-gray-900 transition-colors dark:hover:text-white"
      >
        <ArrowLeft className="h-4 w-4" />
        Back to Dashboard
      </Link>

      <div className="rounded-2xl border border-gray-200 bg-white p-8 shadow-sm dark:border-gray-800 dark:bg-black">
        <div className="mb-8 flex items-center gap-4">
          <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-gray-50 dark:bg-gray-900">
            <Sparkles className="h-6 w-6 text-gray-900 dark:text-white" />
          </div>
          <div>
            <h1 className="text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
              Create New Project
            </h1>
            <p className="text-sm text-gray-500 dark:text-gray-400">
              Deploy your application from a Git repository in seconds.
            </p>
          </div>
        </div>

        {error && (
          <div className="mb-6 flex items-center gap-3 rounded-lg bg-rose-50 p-4 text-sm text-rose-700 dark:bg-rose-500/10 dark:text-rose-400">
            <AlertCircle className="h-5 w-5 shrink-0" />
            <p>{error}</p>
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="space-y-2">
            <label
              htmlFor="name"
              className="text-sm font-medium text-gray-900 dark:text-gray-300"
            >
              Project Name
            </label>
            <input
              id="name"
              type="text"
              required
              placeholder="my-awesome-app"
              className="block w-full rounded-lg border border-gray-200 bg-white px-4 py-3 text-sm transition-all focus:border-gray-900 focus:ring-1 focus:ring-gray-900 dark:border-gray-800 dark:bg-black dark:focus:border-white dark:focus:ring-white"
              value={formData.name}
              onChange={(e) =>
                setFormData({ ...formData, name: e.target.value })
              }
            />
          </div>

          <div className="space-y-2">
            <label
              htmlFor="git_url"
              className="text-sm font-medium text-gray-900 dark:text-gray-300"
            >
              Git Repository URL
            </label>
            <div className="relative">
              <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-4">
                <CodeXml className="h-4 w-4 text-gray-400" />
              </div>
              <input
                id="git_url"
                type="url"
                required
                placeholder="https://github.com/username/repo"
                className="block w-full rounded-lg border border-gray-200 bg-white py-3 pl-11 pr-4 text-sm transition-all focus:border-gray-900 focus:ring-1 focus:ring-gray-900 dark:border-gray-800 dark:bg-black dark:focus:border-white dark:focus:ring-white"
                value={formData.git_url}
                onChange={(e) =>
                  setFormData({ ...formData, git_url: e.target.value })
                }
              />
            </div>
          </div>

          <button
            type="submit"
            disabled={loading}
            className="flex w-full items-center justify-center gap-2 rounded-lg bg-black px-6 py-3.5 text-sm font-medium text-white transition-opacity hover:opacity-90 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-white dark:text-black"
          >
            {loading ? (
              <>
                <Loader2 className="h-4 w-4 animate-spin" />
                Deploying...
              </>
            ) : (
              "Deploy Project"
            )}
          </button>
        </form>
      </div>
    </div>
  );
}
