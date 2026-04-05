# Desafio StatApAnDad 2026

## Introdução

O desafio requer que o código realize o cálculo da média e do desvio padrão de um arquivo com 8,3 bilhões de linhas. A principal dificuldade é fazer esse código executar em um tempo hábil e sem preencher completamente a memória do sistema.

## Primeiras tentativas

Escolhi a linguagem Go pois sabia que seria necessária uma ferramenta que tivesse um desempenho maior e, das linguagens em que tenho proficiência, Go é a mais rápida delas, sendo bem mais eficiente que Python para isso. Em outro momento, posso testar se com bibliotecas como NumPy seria possível fazer esse processamento em Python.

Minha ideia inicial foi calcular somente a contagem (count), a soma e a soma dos quadrados de cada linha, e no final realizar o cálculo da média e desvio padrão. O motivo dessa minha ideia é que realizar uma divisão a cada linha é mais custoso do que uma soma, então somar tudo e calcular somente uma vez no final seria mais eficiente.

Dois problemas foram encontrados:

- Primeiro, depois de pesquisar descobri que essa forma que estava tentando calcular não daria um desvio padrão acurado. Ele é conhecido como um algoritmo ingênuo, pois ignora problemas com o arredondamento de ponto flutuante, especialmente quando as somas resultam em valores muito grandes e o desvio em valores pequenos em comparação, o que causa o acúmulo dos erros de ponto flutuante. Além disso, como eu estava usando um int64, isso levaria a um resultado impreciso.

- Segundo, o programa ainda estava demorando muito para rodar, após testar somente a leitura dos dados sem nenhum cálculo feito, o tempo já era muito elevado, sendo necessário acelerar a forma como o código estava lendo o arquivo.

## Otimizações

Após minhas pesquisas, encontrei um algoritmo melhor para fazer o cálculo: o de Welford. Ele permite calcular a variância de forma incremental, subtraindo a média parcial em cada passo, e fazendo isso em apenas uma passada pelos dados, o que torna o algoritmo ideal para o nosso caso.

```go
count++
delta := x - mean
mean += delta / float64(count)
m2 += delta * (x - mean)
```

E, para os problemas de velocidade de leitura, me recordei de um desafio de código sobre o qual já tinha lido anteriormente, chamado 1BRC (One Billion Row Challenge). O desafio foi iniciado para a linguagem Java, mas rapidamente foi adotado por outras linguagens. Como é bem semelhante ao nosso desafio em alguns aspectos, decidi ler algumas das submissões em Go para me ajudar com essa etapa.

O texto que eu encontrei para me auxiliar nesse processo foi o [One Billion Row Challenge in Golang - From 95s to 1.96s](https://r2p.dev/b/2024-03-18-1brc-go/#:~:text=One%20Billion%20Row%20Challenge%20in%20Golang%20%2D%20From%2095s%20to%201.96s) de Renato Pereira.

A principal otimização mostrada no texto foi a forma como ele saiu de um Scanner básico, padrão do Go, que levou 36s para ler o arquivo dele, para um Scanner otimizado que levou apenas 0,9s para ler o mesmo arquivo.

A forma como ele fez isso foi, inicialmente, adicionando buffers para o bufio (pacote de entradas e saídas do Go com buffer) em diversos tamanhos para encontrar o ideal. O tamanho ideal foi de aproximadamente 4MB, executando o código em cerca de 6,7s. Porém, no texto, o autor nota que o Scanner do bufio tem várias funcionalidades inúteis para o nosso caso, e passa a utilizar a função Read() que também aceita um buffer. Com um buffer do mesmo tamanho, o código executou em 0,9s, sendo a melhor opção para o nosso caso.

```go
buffer := make([]byte, 4<<20) // 4MB mais que isso fica mais lento

for {
    n, err := f.Read(buffer)
    if n > 0 {
        // ...
    }
    if err == io.EOF {
        break
    }
    if err != nil {
        log.Fatal(err)
    }
}
```

Outra otimização que achei útil a partir do texto foi a implementação de um parser customizado. Como sabíamos bem a formatação do texto, que só teria uma casa decimal, pudemos assumir o comportamento da entrada, lendo os caracteres diretamente.

```go
digit := float64(c - '0')
if parsingFrac {
    frac = frac*10 + digit
    fracDiv *= 10
} else {
    num = num*10 + digit
}
```

Com isso, consegui sair de um código que levava algo próximo de uma hora para um código que executa em 15 minutos.

## Feedback do professor

Enviei meus resultados para o professor e fui informado de que a minha média estava errada (com uma diferença de exatamente 2,55), mas o meu desvio padrão estava certo.

Comecei a tentar investigar o código; minha suspeita era de algum erro ou caractere não reconhecido na hora do parser, que estaria sendo adicionado como um número a partir de seu valor ASCII.

Eu estava levando em consideração apenas o . e o \n, sendo ambos os únicos elementos não numéricos que eu acreditava estarem no arquivo. Contudo, após pesquisar na internet sobre erros parecidos, vi que nem sempre a demarcação de quebra de linha é apenas \n; às vezes pode ser \r\n, que é o padrão do Windows.

Após colocar uma condição no código para que, caso encontre esse caractere (\r), ele pule e siga para o próximo, vi que a média chegou ao valor exato.

```go
if c == '\r' {  // corrige o bug que aumentou a media
    continue
}
```
