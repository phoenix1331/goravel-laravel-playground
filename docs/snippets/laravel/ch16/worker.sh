# .env
QUEUE_CONNECTION=redis

# A separate, long-running process (Supervisor, Horizon, or a second container)
php artisan queue:work redis
