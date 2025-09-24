import { z } from 'zod';
import { UserResponseSchema } from '../user/user-response.schema';

// UserResponseSchema is defined in user schemas, just import it

export const SignUpResponseSchema = z.object({
  user: UserResponseSchema,
});

export const SignInResponseSchema = z.object({
  token: z.string(),
  refreshToken: z.string(),
  user: UserResponseSchema,
});

export const SignOutResponseSchema = z.object({
  message: z.string(),
});

export type { UserResponseDto } from '../user/user-response.schema';
export type SignUpResponseDto = z.infer<typeof SignUpResponseSchema>;
export type SignInResponseDto = z.infer<typeof SignInResponseSchema>;
export type SignOutResponseDto = z.infer<typeof SignOutResponseSchema>;
