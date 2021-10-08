drop table bookstore.public.Book;

create table bookstore.public.Book (
  bookId varchar(40),
  bookName varchar(40),
  bookEngName varchar(160),
  authorName varchar(40),
  authorEngName varchar(160),
  category varchar(40),
  note varchar(255),
  isbn varchar(40),
  onStock boolean,
  freeUse boolean,
  imageUrl varchar(255),
  vedioUrl varchar(255),

  PRIMARY KEY (bookId)
);