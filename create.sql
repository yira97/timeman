CREATE TABLE IF NOT EXISTS timeman_todo (
  id           BIGSERIAL                             NOT NULL                         ,
  -- stat: 0: PENDING | 1: WITHDRAW |  2: DOING | 3: CANCEL | 4: TIMEOUT | 4: DONE
  stat         SMALLINT                           NOT NULL     DEFAULT 0           ,
  -- weight > 0 order by created_at DESC
  -- (1, 1001) > (1, 1000) > (-1, 999) > (-1, 1002)
  -- weight < 0 order by created_at ASC
  weight       SMALLINT                           NOT NULL     DEFAULT -1          ,
  avatar       VARCHAR(255)                       NOT NULL     DEFAULT ''          ,
  created_at   TIMESTAMP WITHOUT TIME ZONE        NOT NULL     DEFAULT NOW()       ,
  updated_at   TIMESTAMP WITHOUT TIME ZONE        NOT NULL     DEFAULT NOW()       ,
  withdraw_at  TIMESTAMP WITHOUT TIME ZONE                                         ,
  doing_at     TIMESTAMP WITHOUT TIME ZONE                                         ,
  cancel_at    TIMESTAMP WITHOUT TIME ZONE                                         ,
  dead_time    TIMESTAMP WITHOUT TIME ZONE                                         ,
  done_at      TIMESTAMP WITHOUT TIME ZONE                                         ,
  note         TEXT                                                                ,
  name         VARCHAR(63)                        NOT NULL                         ,
  map_id       BIGINT                                NOT NULL                         ,
  PRIMARY KEY(id)
);

CREATE INDEX todo_list_rank_up ON timeman_todo (weight DESC, created_at DESC);
CREATE INDEX todo_list_rank_bottom ON timeman_todo (weight DESC, created_at ASC);

CREATE TABLE timeline_map (
  id           BIGSERIAL                             NOT NULL                         ,
  name         VARCHAR(63)                           NOT NULL                         ,
  PRIMARY KEY(id)
);