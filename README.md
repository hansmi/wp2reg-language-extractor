# wp2reg-language-extractor

Parse and extract translation strings from the wp2reg firmware used by
Luxtronik 2.x heat pump controllers manufactured and deployed by Alpha
Innotec, NIBE, Novelan and possibly other companies and/or brands.

Tested using wp2reg version 3.85.6.


## Extract firmware

The firmware is packaged as a [tarball][tar] and can be unpacked using GNU tar
and other implementations, e.g.:

```
version=3.85.6 && \
mkdir "wp2reg-${version}" && \
tar --wildcards --to-stdout -xzf "wp2reg-V${version}" "home.wp2reg-V${version}_*" | \
tar -C "wp2reg-${version}" --wildcards -xvf - 'lang_*'
```


## Extract translation strings

```
./wp2reg-language-extractor "wp2reg-${version}"/lang_* > i18n.csv
```

The resulting [CSV file][csv] is encoded using UTF-8 and can be opened using
any spreadsheet software, e.g. [LibreOffice Calc][libreoffice].

[csv]: https://en.wikipedia.org/wiki/Comma-separated_values
[tar]: https://en.wikipedia.org/wiki/Tar_(computing)
[libreoffice]: https://www.libreoffice.org/
