
sqlite3 ./database.db << EOF

DROP TABLE WM_ITEMS;

CREATE TABLE WM_ITEMS(
ID TEXT PRIMARY KEY NOT NULL,
NAME TEXT NOT NULL,
URL_NAME TEXT NOT NULL
);

EOF