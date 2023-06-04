create or replace function gen_chat_service_dataset(
    in_chats_count int, -- Количество генерируемых чатов.
    in_problems_per_chat int, -- Сколько в каждом чате будет проблем.
    in_messages_per_problem int, -- Сколько сообщений в каждой проблеме.
    in_max_body_size int -- Максимальный размер сообщения.
) returns void as
$$
declare
    -- Переменные чатов.
    l_chat_id                    uuid;
    l_client_id                  uuid;

    -- Переменные проблем.
    l_problem_id                 uuid;
    l_manager_id                 uuid;
    l_problem_resolved_at        date;
    l_problem_resolve_request_id uuid;

    -- Переменные сообщений.
    l_author_id                  uuid;
begin
    for i in 1..in_chats_count
        loop
            l_chat_id := gen_random_uuid();
            l_client_id := gen_random_uuid();

            -- Добавляем запись о чате.
            insert into chats(id, client_id, created_at)
            values (l_chat_id, l_client_id, now());

            -- Добавляем проблемы к чату.
            for j in 1..in_problems_per_chat
                loop
                    l_problem_id := gen_random_uuid();
                    l_manager_id := gen_random_uuid();
                    l_problem_resolved_at := now() - '90 days'::interval * random();
                    l_problem_resolve_request_id := gen_random_uuid();

                    -- Только последняя проблема – открытая.
                    if j = in_problems_per_chat then
                        l_problem_resolved_at := null;
                        l_problem_resolve_request_id := null;
                    end if;

                    -- Добавляем запись о проблеме.
                    insert into problems(id, chat_id, manager_id, resolved_at, resolve_request_id, created_at)
                    values (l_problem_id, l_chat_id, l_manager_id, l_problem_resolved_at,
                            l_problem_resolve_request_id, (now() - interval '91 days'));

                    -- Добавляем сообщения к проблеме.
                    for k in 1..in_messages_per_problem
                        loop
                            -- Автор сообщения попеременно то клиент, то менеджер.
                            l_author_id := l_client_id;
                            if k % 2 = 0 then
                                l_author_id := l_manager_id;
                            end if;

                            insert into messages (id, chat_id, problem_id, author_id,
                                                  body, is_visible_for_client, is_visible_for_manager,
                                                  checked_at, initial_request_id, created_at)
                            values (gen_random_uuid(), l_chat_id, l_problem_id, l_author_id,
                                    gen_random_string(cast(ceil(random() * in_max_body_size) as integer)),
                                    true, true, now(), gen_random_uuid(), now());
                        end loop;
                end loop;
        end loop;
end
$$ language plpgsql strict;

create or replace function gen_random_string(length integer) returns text as
$$
declare
    l_chars  text[] := '{a,b,c,d,e,f,g,h,i,j,k,l,m,n,o,p,q,r,s,t,u,v,w,x,y,z,' ||
                       'а,б,в,г,д,е,ё,ж,з,и,й,к,л,м,н,о,п,р,с,т,у,ф,х,ц,ч,ш,щ,ъ,ы,ь,э,ю,я,' ||
                       '0,1,2,3,4,5,6,7,8,9,!,?,.}';
    l_result text   := '';
begin
    for i in 1..length
        loop
            l_result := l_result || l_chars[1 + random() * (array_length(l_chars, 1) - 1)];
        end loop;
    return l_result;
end;
$$ language plpgsql;

select gen_chat_service_dataset(100000, 10, 2, 3000);