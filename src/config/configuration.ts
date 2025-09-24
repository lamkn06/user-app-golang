export default () => ({
  port: parseInt(process.env.PORT ?? '8080', 10),
  apiVersion: 'v1',
  databaseUrl: process.env.DATABASE_URL || 'postgresql://localhost:27017/myapp',
});
