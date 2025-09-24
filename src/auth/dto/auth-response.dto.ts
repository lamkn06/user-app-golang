import { ApiProperty } from '@nestjs/swagger';

export class SignUpResponseDto {
  @ApiProperty()
  user: {
    id: string;
    name: string;
    email: string;
  };
}

export class SignInResponseDto {
  @ApiProperty()
  token: string;

  @ApiProperty()
  refreshToken: string;

  @ApiProperty()
  user: {
    id: string;
    name: string;
    email: string;
  };
}

export class SignOutResponseDto {
  @ApiProperty()
  message: string;
}
