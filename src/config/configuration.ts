interface ConfigModule {
  port: number;
  apiVersion: string;
  databaseUrl: string;
  jwtSecret: string;
}

export default () => ({
  port: parseInt(process.env.PORT ?? '8080', 10),
  apiVersion: 'v1',
  databaseUrl: process.env.DATABASE_URL || 'postgresql://localhost:27017/myapp',
  jwtSecret: process.env.JWT_SECRET || 'default-secret',
});
