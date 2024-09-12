export interface User {
    id: number,
    username: string,
    email: string,
    plan: {
        id: number,
        max_storage: number
    },
    usage: number,
}
