import { redirect } from "next/navigation";
import { requireAdmin } from "@/lib/auth";
import { AdminShell } from "@/components/admin/AdminShell";
export default async function AdminLayout({children}:{children:React.ReactNode}){try{await requireAdmin()}catch{redirect('/login')}return <AdminShell>{children}</AdminShell>}
