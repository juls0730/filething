export interface User {
    id: string,
    username: string,
    email: string,
    plan: {
        id: number,
        max_storage: number
    },
    usage: number,
    created_at: string,
    is_admin: boolean,
}
