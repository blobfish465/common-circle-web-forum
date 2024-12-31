export interface Comment {
    id: number;
    content: string;
    created_at: string;
    updated_at?: string; 
    user_id: number;
    thread_id: number;
}