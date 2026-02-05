export type AdminRole = 'super_admin' | 'ops_admin' | 'finance_admin';
export type AdminStatus = 'active' | 'suspended' | 'inactive';

export interface Admin {
    id: string;
    username: string;
    email: string;
    full_name: string;
    role: AdminRole;
    status: AdminStatus;
    last_login_at?: string;
    created_at: string;
}

export interface AdminLoginRequest {
    username: string;
    password: string;
}

export interface AdminLoginResponse {
    admin_id: string;
    username: string;
    full_name: string;
    role: AdminRole;
    token: string;
}

export interface AuditLog {
    id: string;
    admin_id: string;
    action: string;
    resource_type?: string;
    resource_id?: string;
    description: string;
    ip_address?: string;
    created_at: string;
}