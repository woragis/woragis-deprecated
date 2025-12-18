#!/usr/bin/env python3
"""
Send Test Failure Notification Script

Publishes email notification messages to RabbitMQ for the email-worker to process.
"""

import argparse
import json
import os
import sys
from datetime import datetime

try:
    import pika
except ImportError:
    print("Error: pika library not installed. Run: pip install pika")
    sys.exit(1)


def send_notification(
    rabbitmq_url: str,
    workflow_name: str,
    status: str,
    workflow_url: str,
    commit_sha: str,
    actor: str,
    recipient_email: str,
):
    """Send test failure notification via RabbitMQ."""
    try:
        # Connect to RabbitMQ
        connection = pika.BlockingConnection(pika.URLParameters(rabbitmq_url))
        channel = connection.channel()

        # Declare exchange (should match email-worker configuration)
        exchange = os.getenv("EMAIL_EXCHANGE", "woragis.notifications")
        routing_key = os.getenv("EMAIL_ROUTING_KEY", "emails.send")

        channel.exchange_declare(exchange=exchange, exchange_type="direct", durable=True)

        # Create email message
        subject = f"⚠️ Test Failure: {workflow_name}"
        
        # Format email body
        status_emoji = "❌" if status == "failure" else "⚠️"
        body_html = f"""
        <html>
        <body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
            <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
                <h2 style="color: #d32f2f;">{status_emoji} Test Failure Alert</h2>
                
                <div style="background-color: #fff3cd; border-left: 4px solid #ffc107; padding: 15px; margin: 20px 0;">
                    <strong>Workflow:</strong> {workflow_name}<br>
                    <strong>Status:</strong> {status.upper()}<br>
                    <strong>Triggered by:</strong> {actor}<br>
                    <strong>Time:</strong> {datetime.utcnow().strftime('%Y-%m-%d %H:%M:%S UTC')}
                </div>
                
                <div style="margin: 20px 0;">
                    <h3>Details</h3>
                    <p><strong>Commit:</strong> <code>{commit_sha[:8]}</code></p>
                    <p><strong>Workflow URL:</strong> <a href="{workflow_url}">{workflow_url}</a></p>
                </div>
                
                <div style="margin: 20px 0; padding: 15px; background-color: #f5f5f5; border-radius: 5px;">
                    <p><strong>Action Required:</strong></p>
                    <ul>
                        <li>Review the workflow logs</li>
                        <li>Check for test failures or errors</li>
                        <li>Fix issues and re-run tests</li>
                    </ul>
                </div>
                
                <div style="margin-top: 30px; padding-top: 20px; border-top: 1px solid #ddd; font-size: 12px; color: #666;">
                    <p>This is an automated notification from the CI/CD pipeline.</p>
                </div>
            </div>
        </body>
        </html>
        """
        
        body_text = f"""
Test Failure Alert
==================

Workflow: {workflow_name}
Status: {status.upper()}
Triggered by: {actor}
Time: {datetime.utcnow().strftime('%Y-%m-%d %H:%M:%S UTC')}

Commit: {commit_sha[:8]}
Workflow URL: {workflow_url}

Action Required:
- Review the workflow logs
- Check for test failures or errors
- Fix issues and re-run tests

This is an automated notification from the CI/CD pipeline.
        """

        # Create email envelope (matches email-worker format)
        email_envelope = {
            "user_id": "ci-cd-system",
            "subject": subject,
            "text_message": body_text.strip(),
            "html_message": body_html.strip(),
            "destination": recipient_email,
        }

        # Publish message
        channel.basic_publish(
            exchange=exchange,
            routing_key=routing_key,
            body=json.dumps(email_envelope),
            properties=pika.BasicProperties(
                delivery_mode=2,  # Make message persistent
                content_type="application/json",
            ),
        )

        print(f"✅ Notification sent to {recipient_email} via RabbitMQ")
        print(f"   Exchange: {exchange}")
        print(f"   Routing Key: {routing_key}")

        connection.close()
        return True

    except Exception as e:
        print(f"❌ Failed to send notification: {e}")
        import traceback
        traceback.print_exc()
        return False


def main():
    """Main entry point."""
    parser = argparse.ArgumentParser(description="Send test failure notification")
    parser.add_argument("--workflow", required=True, help="Workflow name")
    parser.add_argument("--status", required=True, help="Workflow status")
    parser.add_argument("--url", required=True, help="Workflow URL")
    parser.add_argument("--commit", required=True, help="Commit SHA")
    parser.add_argument("--actor", required=True, help="Actor who triggered workflow")
    parser.add_argument("--email", required=True, help="Recipient email address")
    parser.add_argument(
        "--rabbitmq-url",
        default=os.getenv("RABBITMQ_URL", "amqp://test:test@localhost:5673/test"),
        help="RabbitMQ connection URL",
    )

    args = parser.parse_args()

    success = send_notification(
        rabbitmq_url=args.rabbitmq_url,
        workflow_name=args.workflow,
        status=args.status,
        workflow_url=args.url,
        commit_sha=args.commit,
        actor=args.actor,
        recipient_email=args.email,
    )

    sys.exit(0 if success else 1)


if __name__ == "__main__":
    main()
