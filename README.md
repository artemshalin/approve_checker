# approve checker

Утилита позволяет подключить `approve rules` в бесплатной версии GitLab CE.

## Approve rules для merge request в бесплатной версии GitLab CE

### Шаг 1: Включаем проверку слияния

Settings → Merge requests. Для "Merge checks" ставим галочку "Pipelines must succeed".

### Шаг 2: Настраиваем "Protected branches"

Переходим Settings → Repository → Protected branches и для необходимых веток включаем защиту. Например, это могут быть ветки: main и stage.

### Шаг 3: Настройка доступа к репозиторию

Для работы утилиты необходимо выполнить следующие условия:

- Пользователь, под которым работает утилита, должен быть в проекте с правами `Reporter` или выше.
- Указать в переменных проекта переменную окружения `GITLAB_TOKEN` токен пользователя .

Можно использовать свой PAT или создать нового пользователя, получить его PAT, затем добавить его в проект с правами `Reporter` или выше.

Шаги следующие:

1. Убеждаемся, что пользователь есть в проекте. Добавляем если, его нет.
2. Создаем для этого пользователя персональный токен с доступом к API GitLab (PAT): [Мануал по созданию PAT](https://docs.gitlab.com/user/profile/personal_access_tokens/#create-a-personal-access-token)
3. Переходим Settings → CI/CD → Variables → Key `GITLAB_TOKEN` → Value `Ваш PAT` → Ставим галочку Masked → Ставим или убираем галочку Protected → Ok

### Шаг 4: Настройка переменных окружения

Доступны следующие настройки:

- `GITLAB_TOKEN` - Токен доступа пользователя, который есть в проекте. Под этим пльзователем будет работать утилита.
- `APPROVE_MIN_APPROVAL_ROLE` - Минимальная роль, одобрение которой учитывает утилита. Возможны следующие значения:
  - MinimalAccess = 5
  - Guest         = 10
  - Planner       = 15
  - Reporter      = 20
  - Developer     = 30
  - Maintainer    = 40
  - Owner         = 50
  - Admin         = 60
- `APPROVE_APPROVAL_AUTHORS` - Персональные имена учетных записей, одобрения которых учитываются при согласовании. Если учетных записей несколько, то необходимо разделить их запятой без пробелов.
- `APPROVE_MIN_APPROVAL_COUNT` - Минимальное количество одобрений, которое необходимо получить, для согласования.

Для настройки переходим Settings → CI/CD → Variables

### Шаг 5: Подключение компонента к своему репозиторию

Чтобы подключить компонент к своему репозиторию необходимо в файле `.gitlab-ci.yml` добавить следующий кусок yml:

```yml
# stages может быть любой
stages:
- test

check-mr-approve:
  stage: test    
  image: 
    name: ghcr.io/artemshalin/approve_checker:v1.0.1
    entrypoint: [""]
  tags: 
    - docker
  script:
    - /usr/local/bin/approve_checker
  rules:
    - if: '$CI_MERGE_REQUEST_TARGET_BRANCH_NAME == "main" || $CI_MERGE_REQUEST_TARGET_BRANCH_NAME == "stage" || $CI_MERGE_REQUEST_TARGET_BRANCH_NAME == "dev"'
```
