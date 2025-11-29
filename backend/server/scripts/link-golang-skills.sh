#!/bin/bash

# Link Golang skill to all projects for the user
USER_ID="6ad0d828-f605-45fc-a545-3441e17a015c"

echo "Linking Golang skill to projects..."

# Get Golang skill ID
GOLANG_SKILL_ID=$(docker exec woragis-database psql -U postgres -d woragis -t -c "SELECT id FROM skills WHERE name = 'Golang' LIMIT 1;" | tr -d ' ')

if [ -z "$GOLANG_SKILL_ID" ]; then
    echo "Error: Golang skill not found"
    exit 1
fi

echo "Golang skill ID: $GOLANG_SKILL_ID"

# Get all project IDs for the user
PROJECT_IDS=$(docker exec woragis-database psql -U postgres -d woragis -t -c "SELECT id FROM projects WHERE user_id = '$USER_ID';" | tr -d ' ')

count=0
for PROJECT_ID in $PROJECT_IDS; do
    if [ ! -z "$PROJECT_ID" ]; then
        docker exec woragis-database psql -U postgres -d woragis -q -c "
            INSERT INTO project_skills (project_id, skill_id, created_at)
            VALUES ('$PROJECT_ID', '$GOLANG_SKILL_ID', NOW())
            ON CONFLICT (project_id, skill_id) DO NOTHING;
        " > /dev/null 2>&1
        count=$((count + 1))
    fi
done

echo "Linked Golang skill to $count projects"

# Verify
docker exec woragis-database psql -U postgres -d woragis -c "
    SELECT p.name, COUNT(ps.skill_id) as skills_count
    FROM projects p
    LEFT JOIN project_skills ps ON ps.project_id = p.id
    WHERE p.user_id = '$USER_ID'
    GROUP BY p.id, p.name
    ORDER BY p.created_at DESC
    LIMIT 5;
" -t

