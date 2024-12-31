export interface Thread {
    id: number;
    user_id: number; 
    title: string;
    content: string;
    created_at: string; 
    updated_at?: string; 
    category_id: number;
}