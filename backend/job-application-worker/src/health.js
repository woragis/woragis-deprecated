/**
 * Health check module for Job Application Worker.
 * Checks service availability and RabbitMQ connection.
 */

let rabbitmqConnection = null;

/**
 * Set RabbitMQ connection for health checks
 */
export const setRabbitMQConnection = (conn) => {
  rabbitmqConnection = conn;
};

/**
 * Check health status
 */
export const checkHealth = () => {
  const checks = [];
  
  // Service check (always ok if endpoint is reachable)
  checks.push({
    name: 'service',
    status: 'ok'
  });
  
  // RabbitMQ check
  const rabbitmqCheck = checkRabbitMQ();
  checks.push(rabbitmqCheck);
  
  // Determine overall status
  const hasErrors = checks.some(check => check.status === 'error');
  const status = hasErrors ? 'unhealthy' : 'healthy';
  
  return {
    status,
    checks
  };
};

/**
 * Check RabbitMQ connection
 */
const checkRabbitMQ = () => {
  if (!rabbitmqConnection) {
    return {
      name: 'rabbitmq',
      status: 'error',
      message: 'not configured'
    };
  }
  
  try {
    // Check if connection is closed (amqplib specific)
    // amqplib connections don't have readyState, check if connection exists
    if (!rabbitmqConnection.connection) {
      return {
        name: 'rabbitmq',
        status: 'error',
        message: 'connection not initialized'
      };
    }
    
    // For amqplib, we can't easily check if connection is open without a ping
    // Assume ok if connection object exists
    return {
      name: 'rabbitmq',
      status: 'ok'
    };
  } catch (error) {
    return {
      name: 'rabbitmq',
      status: 'error',
      message: error.message
    };
  }
};
