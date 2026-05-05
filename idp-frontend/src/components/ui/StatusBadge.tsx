import { cn } from "@/lib/utils";
import { ProjectStatus } from "@/types";

interface StatusBadgeProps {
  status: ProjectStatus;
  className?: string;
}

export function StatusBadge({ status, className }: StatusBadgeProps) {
  const getStatusConfig = (status: ProjectStatus) => {
    switch (status) {
      case "READY":
        return {
          dot: "bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.6)]",
          text: "text-emerald-700 dark:text-emerald-400",
          bg: "bg-emerald-50 dark:bg-emerald-500/10",
          label: "Ready",
        };
      case "PENDING":
      case "BUILDING":
        return {
          dot: "bg-blue-500 animate-pulse shadow-[0_0_8px_rgba(59,130,246,0.6)]",
          text: "text-blue-700 dark:text-blue-400",
          bg: "bg-blue-50 dark:bg-blue-500/10",
          label: status === "BUILDING" ? "Building" : "Pending",
        };
      case "FAILED":
        return {
          dot: "bg-rose-500 shadow-[0_0_8px_rgba(244,63,94,0.6)]",
          text: "text-rose-700 dark:text-rose-400",
          bg: "bg-rose-50 dark:bg-rose-500/10",
          label: "Failed",
        };
      default:
        return {
          dot: "bg-gray-400",
          text: "text-gray-700 dark:text-gray-400",
          bg: "bg-gray-50 dark:bg-gray-500/10",
          label: status,
        };
    }
  };

  const config = getStatusConfig(status);

  return (
    <div
      className={cn(
        "inline-flex items-center gap-2 px-2.5 py-0.5 rounded-full text-xs font-medium border transition-all duration-300",
        config.bg,
        config.text,
        "border-transparent",
        className
      )}
    >
      <span className={cn("h-2 w-2 rounded-full", config.dot)} />
      {config.label}
    </div>
  );
}
