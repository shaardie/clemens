# TODO: Bugs in der Transposition Table

**Status:**
- 🔴 Bug 4 – offen (Matt-Score in der TT)
- 🔴 Bug 5 – offen (TT: Duplikate derselben Stellung)
- 🔴 Bug 6 – offen (TT: `setAge`-Overflow)
- 🔴 Bug 7 – offen (TT: `Reset()` setzt `hashEntries` nicht zurück)
- 🔴 Bug 8 – offen (TT: `ucinewgame` setzt die TT nicht zurück)

## Bug 4: Matt-Score im TT wird nur beim Lesen, nicht beim Speichern angepasst

**Datei:** `pkg/search/transpositiontable/transposition_table.go:69-74`

Beim Speichern (`search.go:400`) wird `bestScore` roh abgelegt. Für einen
Matt-Zug ist das `INF - childPly` (also **relativ zum Wurzel-Ply**):

- Kind-Knoten (Ply `P`) ist matt → liefert `-INF + P`
- Eltern-Knoten speichert `score = -(-INF + P) = INF - P`

Beim Lesen in `Get` wird dann das **aktuelle** `ply` abgezogen:

```go
if score > evaluation.INF-100 {
    score -= int16(ply)
} else if score < -evaluation.INF+100 {
    score += int16(ply)
}
```

Da der gespeicherte Wert das **gespeicherte** `P` bereits enthält, ergibt sich
in einer transponierten Stellung z. B. `INF - 11 - 5 = INF - 16` statt
`INF - 6`.

**Wirkung:** Verfälschte Mattscores → falsche Cutoffs / Move-Wahl in
Mattstellungen bei Transpositionen.

**Fix-Muster:** Beide Seiten anpassen, d. h. node-relativ speichern:
- Beim Speichern: `score += ply` für positive Matts, `score -= ply` für negative.
- Beim Lesen: `score -= ply` für positive, `score += ply` für negative.

Dafür muss `PotentiallySave` das `ply` übergeben bekommen (aktuell nicht).

---

## Bug 5: `PotentiallySave` kann Duplikate derselben Stellung anlegen

**Datei:** `pkg/search/transpositiontable/transposition_table.go:97-112`

Die Ersetzungsschleife sucht nur nach leeren/schlechteren Slots und prüft
**nicht**, ob dieselbe `zobristHash` schon im Bucket liegt. Liegt ein
tieferer/neuerer Eintrag weiter hinten und vorne ist ein freier Slot, wird ein
**schlechteres Duplikat** vorne eingefügt. `Get` scannt von Index 0 und nimmt
den **ersten** Treffer → es kann das schlechtere Duplikat verwenden.

**Fix:** Beim Speichern zuerst nach derselben Hash suchen und diesen Eintrag
aktualisieren.

---

## Bug 6: `setAge`-Overflow für `HalfMoveClock >= 64`

**Datei:** `pkg/search/transpositiontable/ttentry.go:26`

`age` ist der `HalfMoveClock` (uint8). Bei `age >= 64` läuft `age << 2`
über (`64<<2 = 256 → 0`). Solche Einträge erscheinen als „ganz frisch"
(Age 0) und werden seltener ersetzt. Betrifft nur die Ersetzungsheuristik.

---

## Bug 7: `Reset()` setzt `hashEntries` nicht zurück

**Datei:** `pkg/search/transpositiontable/transposition_table.go:39-41`

`clear(tt[:])` leert die Tabelle, aber `hashEntries` bleibt stehen.
`HashFull()` meldet danach zu hohe Werte. Aktuell egal, weil
`transpositiontable.Reset()` in `pkg/uci/input.go:31` **auskommentiert** ist.

---

## Bug 8: `ucinewgame` setzt die TT nicht zurück

**Datei:** `pkg/uci/input.go:29-32`

Bei `ucinewgame` wird nur `game.New()` ausgeführt, die TT bleibt bestehen.
Da die Hashes jetzt konsistent sind, sind Einträge aus alten Partien technisch
gültig, aber die Age-Heuristik (Halfmove-Clock) unterscheidet keine Partien.
