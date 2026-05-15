import type { ButtonHTMLAttributes, AnchorHTMLAttributes, ReactNode } from "react";
import Link from "next/link";

const base = "inline-flex items-center justify-center rounded-md px-5 py-3 text-sm font-bold transition focus:outline-none focus:ring-2 focus:ring-red-500 disabled:opacity-50";
const variants = { primary: "bg-red-600 text-white hover:bg-red-700", secondary: "bg-white/10 text-white hover:bg-white/20", ghost: "text-white hover:bg-white/10" };

export function Button({ variant = "primary", className = "", ...props }: ButtonHTMLAttributes<HTMLButtonElement> & { variant?: keyof typeof variants }) {
  return <button className={`${base} ${variants[variant]} ${className}`} {...props} />;
}
export function ButtonLink({ href, children, variant = "primary", className = "", ...props }: AnchorHTMLAttributes<HTMLAnchorElement> & { href: string; children: ReactNode; variant?: keyof typeof variants }) {
  return <Link href={href} className={`${base} ${variants[variant]} ${className}`} {...props}>{children}</Link>;
}
