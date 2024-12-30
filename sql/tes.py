# sql/tes.py
import sqlite3
con = sqlite3.connect('movies.db')
cur = con.cursor()  
cur.execute("CREATE TABLE  movies ( title TEXT, genre INTEGER, year INTEGER)")
movies = [
    ("The Shawshank Redemption", 1994, 9.2),
    ("The Godfather", 1972, 9.2),
    ("The Dark Knight", 2008, 9.0),
    ("12 Angry Men", 1957, 9.0),
]
cur.executemany("INSERT INTO movies  VALUES (?,?,?)",movies)
data = cur.execute( " SELECT rowid, * FROM movies")
for row in data:
    print(row)

con.commit()
con.close()