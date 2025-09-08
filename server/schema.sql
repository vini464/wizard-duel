DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS cards;
DROP TABLE IF EXISTS user_cards;
DROP TABLE IF EXISTS card_stock;


CREATE TABLE cards (
  name TEXT NOT NULL PRIMARY KEY,
  type TEXT NOT NULL,
  rarity TEXT NOT NULL,
  cost INTEGER NOT NULL
  effect1 TEXT NOT NULL,
  amount1 INTEGER NOT NULL,
  effect2 TEXT,
  amount2 INTEGER,
);


CREATE TABLE users (
  username TEXT NOT NULL,
  password TEXT NOT NULL,
  coins TEXT NOT NULL,
);


CREATE TABLE user_cards (
  username TEXT NOT NULL REFERENCES users(username) ON DELETE CASCADE,
  cardname TEXT NOT NULL REFERENCES cards(cardname) ON DELETE CASCADE,
  quantity INTEGER DEFAULT 1 CHECK (quantity>=0),
  PRIMARY KEY (username, cardname),
);


CREATE TABLE card_stock (
  cardname TEXT NOT NULL REFERENCES cards(cardname) ON DELETE CASCADE,
  quantity INTEGER NOT NULL CHECK (quantity>=0)
  PRIMARY KEY (cardname)
);
