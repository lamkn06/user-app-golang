import { z } from 'zod';

export const ListRequestSchema = z.object({
  page: z.coerce.number().min(1).default(1),
  limit: z.coerce.number().min(1).default(10),
});

export type ListRequestDto = z.infer<typeof ListRequestSchema>;
