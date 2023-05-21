# DNS Resolver

This project is based on https://implement-dns.wizardzines.com/index.html (it is a Jupyter notebook, the zip file is available in the source folder in case the website gets deleted !)

The purpose of this project is to allow me to practise writing programs in Golang and to help me understand a bit better how DNS works !

## Some useful documentation :
- About DNS itself :
  https://medium.com/@openmohan/dns-basics-and-building-simple-dns-server-in-go-6cb8e1cfe461
- DNS Headers :
http://www.tcpipguide.com/free/t_DNSMessageHeaderandQuestionSectionFormat.htm
- DNS Questions :
https://www.zytrax.com/books/dns/ch15/#question
- Other stuff on DNS Packet structure :
https://mislove.org/teaching/cs4700/spring11/handouts/project1-primer.pdf
- On compression :
http://www.tcpipguide.com/free/t_DNSNameNotationandMessageCompressionTechnique-2.htm
- Implementation made by Google, which was useful to see what I did wrong !
https://github.com/google/gopacket/blob/master/layers/dns.go

## Run the project

To run, use the following command and whitelist the dns-resolver folder in your favorite antivirus :
```
go build -o .\out\dns-resolver && .\out\dns-resolver
```

## Next step

Well I do not really plan on doing more but this project could use :
- More error handling
- Write some unit tests on the decoding / encoding !
- Protection against loop for DNS compression (this could allow some attack)
- Implement some DNS caching ?
- Handle more type of Records ?
- Make a CLI ?