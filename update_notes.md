# Opis zmian oraz odpowiedzi na pytania

## Zmiany

### Baza danych

Do bazy danych zostały dodane twarde ograniczenia UNIQUE które zapewniają, że w bazie nie mogą pojawić się duplikaty słów, tłumaczeń oraz przykładów

## Resolver

W nowej wersji mutacji wykożystywany jest gorm Transaction, który obsługuje autoamtycznie tranzakcje z bazą danych. Dodatkowo mutacje zwracają uwagę czy błąd podczas operacji interakcji z bazą danych jest przewidywany, czy jest to jakiś inny błąd.


## Odpowiedzi na pytania

1. Jaki jest cel dodatkowej konfiguracji bazy danych (SetMaxOpenConns itd.)?
2. Co daje użycie SET CONSTRAINTS ALL IMMEDIATE;? I czemu jest zdefiniowane tylko w testach?

1) Był to fragment kodu,który służył do sprawdzenia jakie wyniki powinny dać testy dotyczące współbierzności. Na wczesnym etapie projektowania testów nie tworzyłem testowego serwera, więc nie współbieżność nie była dobrze obsłużona, po ustawieniu SetMaxOpenConns(1) byłem w stanie zasymulować działanie współbierzne. Po zmianie podejścia do testów fragment ten został przypadkiem, powinien zostać zakomentowany

2) Chciałem sprawić, żeby testy pokazywały również, że wszelkie ograniczenia sprawdzane były od razu po zapytaniu, a nie tylko na końcu transakcji. Jednak nie wpływa to bezpośrednio na wynik testu, więc również aby utrzywmać zbliżoność testu do src, linia ta również została usunięta.



## Testy

Bardzo prosił bym, żeby za pierwszym razem po odpaleniu bazy danych testowej nie używać flagi -count, ponieważ kontener czasami ma problemy z połączeniem. Przy kojenych odpaleniach, już jak można to przetestować. (Sam odpaliłem testy 10 razy z flagą -count 8, wszytkie przeszły)