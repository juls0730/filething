export interface File {
    name: string,
    is_dir: boolean,
    size: number,
    last_modified: string,
    toggled: "checked" | "some" | "unchecked",
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
    } | {},
    status: {
        error: boolean,
        aborted: boolean,
        code: number,
        message: string
    } | {},
}

export interface UploadResponse {
    usage: number,
    file: File
}