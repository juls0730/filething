export interface User {
    id: string,
    username: string,
    email: string,
    plan: Plan,
    usage: number,
    created_at: string,
    is_admin: boolean,
}

export interface Plan {
    id: number,
    max_storage: number
}