import { z } from 'zod';

export const SignUpSchema = z.object({
  email: z.string().email('Invalid email format'),
  password: z.string().min(6, 'Password must be at least 6 characters'),
});

export type SignUpDto = z.infer<typeof SignUpSchema>;
