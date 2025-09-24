import { z } from 'zod';

export const UserResponseSchema = z.object({
  id: z.string(),
  name: z.string(),
  email: z.string().email(),
});

export type UserResponseDto = z.infer<typeof UserResponseSchema>;
