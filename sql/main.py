# sql/main.py
import sqlite3
con = sqlite3.connect("movies.db")

con.close()