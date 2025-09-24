import { z } from 'zod';

export const PaginationSchema = z.object({
  page: z.coerce.number().min(1).default(1),
  limit: z.coerce.number().min(1).max(100).default(10),
  sortBy: z.string().optional(),
  sortOrder: z.enum(['asc', 'desc']).default('desc'),
});

export type PaginationDto = z.infer<typeof PaginationSchema>;

export type ListResponseDto<T = any> = {
  data: T[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
};
