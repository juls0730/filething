export interface User {
    id: number,
    username: string,
    email: string,
    plan: {
        id: number,
        max_storage: number
    }
}

export interface FileUpload {
    id: string,
    uploading: boolean,
    file: File,
    startTime: number,
    speed: number,
    remainingTime: number,
    controller: XMLHttpRequest,
    length: {
        total: number,
        loaded: number,
    },
    status: {
        error: boolean,
        aborted: boolean,
        code: number,
        message: string
    },
}