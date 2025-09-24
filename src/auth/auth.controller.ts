import {
  Controller,
  Post,
  Body,
  UseGuards,
  Request,
  UsePipes,
} from '@nestjs/common';
import {
  ApiTags,
  ApiOperation,
  ApiResponse,
  ApiBearerAuth,
  ApiProperty,
} from '@nestjs/swagger';
import { AuthService } from './auth.service';
import { SignUpSchema, SignInSchema } from '../schemas';
import type { SignUpDto, SignInDto } from '../schemas';
import {
  SignUpResponseDto,
  SignInResponseDto,
  SignOutResponseDto,
} from '../dto';
import { JwtAuthGuard } from '../common/guards/jwt-auth.guard';
import { ZodValidationPipe } from '../common/pipes/zod-validation.pipe';

@ApiTags('auth')
@Controller('auth')
export class AuthController {
  constructor(private readonly authService: AuthService) {}

  @Post('signup')
  @UsePipes(new ZodValidationPipe(SignUpSchema))
  @ApiOperation({ summary: 'Sign up a new user' })
  @ApiResponse({
    status: 201,
    description: 'User created successfully',
    type: SignUpResponseDto,
  })
  @ApiResponse({ status: 409, description: 'User already exists' })
  async signUp(@Body() signUpDto: SignUpDto): Promise<SignUpResponseDto> {
    return this.authService.signUp(signUpDto) as any;
  }

  @Post('signin')
  @UsePipes(new ZodValidationPipe(SignInSchema))
  @ApiOperation({ summary: 'Sign in a user' })
  @ApiResponse({
    status: 200,
    description: 'User signed in successfully',
    type: SignInResponseDto,
  })
  @ApiResponse({ status: 401, description: 'Invalid credentials' })
  async signIn(@Body() signInDto: SignInDto): Promise<SignInResponseDto> {
    return this.authService.signIn(signInDto) as any;
  }

  @Post('signout')
  @UseGuards(JwtAuthGuard)
  @ApiBearerAuth()
  @ApiOperation({ summary: 'Sign out a user' })
  @ApiResponse({
    status: 200,
    description: 'User signed out successfully',
    type: SignOutResponseDto,
  })
  @ApiResponse({ status: 401, description: 'Unauthorized' })
  async signOut(@Request() req): Promise<SignOutResponseDto> {
    return this.authService.signOut();
  }
}
