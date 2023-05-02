-- TODO: REMOVER! Dados de exemplo para testes.

-- usuário
INSERT INTO forum_users (nickname, email, active)
VALUES ('admin', 'example@example.com', true);

-- A tabela de topics é recursiva para permitir topicos aninhados.
-- abaixo alguns exemplos para teste

-- tópicos Golang
INSERT INTO forum_topics (id, title, slug, description, created_by, parent_id)
VALUES (1, 'Golang', 'golang', 'Tópico de Golang', 1, 0);

INSERT INTO forum_topics (id, title, slug, description, created_by, parent_id)
VALUES (2, 'Geral', 'geral', 'Tópico geral', 1, 1);

INSERT INTO forum_topics (id, title, slug, description, created_by, parent_id)
VALUES (3, 'Instalação', 'instalacao', 'Tópico de instalação', 1, 1);

INSERT INTO forum_topics (id, title, slug, description, created_by, parent_id)
VALUES (4, 'Desenvolvimento', 'desenvolvimento', 'Tópico de desenvolvimento', 1, 1);

-- tópicos Assembly

INSERT INTO forum_topics (id, title, slug, description, created_by, parent_id)
VALUES (5, 'Assembly', 'assembly', 'Tópico de Assembly', 1, 0);

INSERT INTO forum_topics (id, title, slug, description, created_by, parent_id)
VALUES (6, 'Geral', 'geral', 'Tópico geral', 1, 5);

INSERT INTO forum_topics (id, title, slug, description, created_by, parent_id)
VALUES (7, 'Instalação', 'instalacao', 'Tópico de instalação', 1, 5);

INSERT INTO forum_topics (id, title, slug, description, created_by, parent_id)
VALUES (8, 'Mercado de trabalho', 'mercado-de-trabalho', 'Tópico de mercado de trabalho', 1, 5);

-- threads
-- trheads contem um conjunto de posts (um no minimo) e um topico, elas contem o
-- titulo/assunto e outros dados relacionados a thread.

INSERT INTO forum_threads (id, topic_id, title, slug, description, created_by)
VALUES (1, 1, 'Como fazer um hello world em Golang', 'hello-world', 'Como fazer um hello world em Golang', 1);

-- posts

INSERT INTO forum_posts (id, thread_id, content, created_by)
VALUES (1, 1, 'Olá, mundo!', 1);

