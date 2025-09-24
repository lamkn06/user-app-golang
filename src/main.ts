import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import {
  NestFastifyApplication,
  FastifyAdapter,
} from '@nestjs/platform-fastify';
import { ConfigService } from '@nestjs/config';

async function bootstrap() {
  const app = await NestFactory.create<NestFastifyApplication>(
    AppModule,
    new FastifyAdapter(),
  );

  const configService: ConfigService = app.get(ConfigService);

  app.setGlobalPrefix(`api/${configService.get<string>('apiVersion')}`);

  const port = configService.get<number>('port') || 8080;
  await app.listen(port, '0.0.0.0');
}
bootstrap();
