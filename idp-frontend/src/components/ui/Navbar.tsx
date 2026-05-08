"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  ChevronRight,
  LayoutGrid,
  Plus,
  User,
  Box,
  Command,
} from "lucide-react";
import { cn } from "@/lib/utils";

export function Navbar() {
  const pathname = usePathname();
  const pathSegments = pathname.split("/").filter(Boolean);

  return (
    <nav className="sticky top-0 z-50 w-full border-b border-gray-200 bg-white/80 backdrop-blur-md dark:border-gray-800 dark:bg-black/80">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="flex h-16 items-center justify-between">
          <div className="flex items-center gap-4">
            <Link href="/dashboard" className="flex items-center gap-2">
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-black text-white dark:bg-white dark:text-black">
                <Command className="h-5 w-5" />
              </div>
              <span className="text-sm font-bold tracking-tight text-gray-900 dark:text-white uppercase">
                Integrated Deployment Platform
              </span>
            </Link>

            <div className="h-6 w-px bg-gray-200 dark:bg-gray-800 mx-2" />

            <div className="flex items-center gap-1.5 text-sm text-gray-500">
              <Link
                href="/dashboard"
                className="hover:text-gray-900 dark:hover:text-white transition-colors"
              >
                Dashboard
              </Link>
              {pathSegments.length > 0 && pathSegments[0] !== "dashboard" && (
                <>
                  <ChevronRight className="h-4 w-4" />
                  <span className="font-medium text-gray-900 dark:text-white capitalize">
                    {pathSegments[0]}
                  </span>
                </>
              )}
            </div>
          </div>

          <div className="flex items-center gap-4">
            <Link
              href="/new"
              className="hidden items-center gap-2 rounded-md bg-black px-3 py-1.5 text-sm font-medium text-white transition-opacity hover:opacity-90 sm:flex dark:bg-white dark:text-black"
            >
              <Plus className="h-4 w-4" />
              New Project
            </Link>
            <div className="h-8 w-8 overflow-hidden rounded-full border border-gray-200 bg-gray-100 dark:border-gray-800 dark:bg-gray-900">
              <User className="h-full w-full p-1.5 text-gray-400" />
            </div>
          </div>
        </div>
      </div>
    </nav>
  );
}
