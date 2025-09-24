import { z } from 'zod';
import { UserResponseSchema } from './user-response.schema';

export const ListResponseSchema = z.object({
  data: z.array(UserResponseSchema),
  total: z.number(),
  page: z.number(),
  limit: z.number(),
  totalPages: z.number(),
});

export const UserListResponseSchema = ListResponseSchema;

export type UserListResponseDto = z.infer<typeof UserListResponseSchema>;
