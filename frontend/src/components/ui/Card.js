import React from "react";
import clsx from "clsx";

export function Card({ children, className = "", ...props }) {
  return (
    <div
      className={clsx(
        "bg-white rounded-lg border border-gray-200 shadow-sm",
        className
      )}
      {...props}
    >
      {children}
    </div>
  );
}

export function CardHeader({ children, className = "", ...props }) {
  return (
    <div
      className={clsx("px-6 py-4 border-b border-gray-200", className)}
      {...props}
    >
      {children}
    </div>
  );
}

export function CardContent({ children, className = "", ...props }) {
  return (
    <div className={clsx("px-6 py-4", className)} {...props}>
      {children}
    </div>
  );
}

export function CardFooter({ children, className = "", ...props }) {
  return (
    <div
      className={clsx(
        "px-6 py-4 border-t border-gray-200 bg-gray-50 rounded-b-lg",
        className
      )}
      {...props}
    >
      {children}
    </div>
  );
}
