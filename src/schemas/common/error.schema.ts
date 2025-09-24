import { z } from 'zod';

export const ErrorDetailSchema = z.object({
  field: z.string(),
  message: z.string(),
});

export const ErrorResponseSchema = z.object({
  message: z.string(),
  errors: z.array(ErrorDetailSchema).optional(),
  statusCode: z.number(),
});

export type ErrorDetail = z.infer<typeof ErrorDetailSchema>;
export type ErrorResponse = z.infer<typeof ErrorResponseSchema>;
