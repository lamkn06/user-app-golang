import { Injectable } from '@nestjs/common';
import { PrismaService } from '../config/prisma/prisma.service';
import {
  CreateUserDto,
  UserResponseDto,
  ListRequestDto,
  UserListResponseDto,
} from '../schemas';

@Injectable()
export class UserService {
  constructor(private readonly prisma: PrismaService) {}

  async createUser(createUserDto: CreateUserDto): Promise<UserResponseDto> {
    const user = await this.prisma.user.create({
      data: {
        ...createUserDto,
        password: 'default-password', // This should be handled properly in real app
      },
    });

    return {
      id: user.id,
      name: user.name || '',
      email: user.email,
    };
  }

  async getUsers(listRequest: ListRequestDto): Promise<UserListResponseDto> {
    const { limit, page } = listRequest;
    const offset = ((page || 1) - 1) * (limit || 10);

    // Get total count
    const total = await this.prisma.user.count();

    // Get paginated users
    const users = await this.prisma.user.findMany({
      skip: offset,
      take: limit,
      orderBy: { createdAt: 'desc' },
      select: {
        id: true,
        name: true,
        email: true,
      },
    });

    // Map users to ensure name is never null
    const mappedUsers = users.map((user) => ({
      id: user.id,
      name: user.name || '',
      email: user.email,
    }));

    return {
      data: mappedUsers,
      total,
      page: page || 1,
      limit: limit || 10,
      totalPages: Math.ceil(total / (limit || 10)),
    };
  }

  async getUserById(id: string): Promise<UserResponseDto> {
    const user = await this.prisma.user.findUniqueOrThrow({
      where: { id },
      select: {
        id: true,
        name: true,
        email: true,
      },
    });

    return {
      id: user.id,
      name: user.name || '',
      email: user.email,
    };
  }
}
