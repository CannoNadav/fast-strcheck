
The task: Given a set of 'dictionary' strings, the task is to build a service that handles high volumes of input strings, efficiently identifying which of them are a permutation of one of the strings in the dictionary. Given a string s as input the service must return all the strings in the dictionary who contain exactly the same letters as s. The service must also provide real-time data on the number of requests served and CPU time spent processing them, under a dedicated endpoint. The challange is to making a retrival in minimum time and minimum memory resources.



For instructions on running the app, skip down to "How to run"



A top down overview of the app:

The application is written in go using the Echo framework and packaged inside a docker container.
I chose this language as it offers high performance and good concurrency support, and the echo framework is mainstream, minimalistic, and high performance meaning it will be reliable and not add a lot of overhead to the application.

How the backend works:
when the server is initialized it accepts a parameter of type WordsCatalogFactory. WordsCatalog is the common interface that the two solution algorithms implement to interoperate inside the server. The type of words catalog the server will be initialized with can be configured by the command line parameter 'WordsCatalogAlg' which has two possible values: MapWordsCatalog and MemBlockWordsCatalog
These represent the two implementations of the main algorithm for the task. more on that later.

The task informally defined a relation similar-words on the set of words in [a-z]*:

similar-words(w1, w2) iff: 
1) w1 is a letter permutation of w2, and w2 is listed in 'words_clean.txt'
2) w1!=w2 (anti-reflexivity)

if we remove the second requirement(anti-reflexivity) then we get an equivalence relation. Lets call it similar-words'. given a working implementation for similar-words'(i.e one that given a word returns its equivalence group) , it is easy to construct an implementation for similar-words. For every input word w1 just find its equivalence group in similar-words' and remove w1 from it, if it is present there. That's basically what both implementations are doing, just in a different way.

Since it would have been inefficient to search through the entire dictionary and collect all equivalent words for every similar words request, in both implementations I preprocess the supplied dictionary to extract its equivalence groups, and store each such group in the respective data structure for that algorithm, for efficient search later.

Let's define a term - "the canonical representation" of a word is the word, where its letters were lexicographically sorted. e.g the canonical representation of the word "tram" is "amrt". Since all words in each equivalence group contain different permutations of the same letters, their canonical representation is the same.

MapWordsCatalog - This was the first implementation I've made. It's a map where each equivalence group in similar-words' is stored as a map value, and the map key is the common canonical representation of the strings in the group. When a similar words request comes in, we sort the string obtained from the request's word parameter, and look for its equivalence group in the map. if such a group is found then we return it, after removing the searched word from the result.

pros: simple both conceptually and to implement. fast since all we need to do when receiving a request is sort an input string and look for the sorted string in a map.
	cons: it wastes a lot of memory. First, there's the map overhead - memory overheads related to saving the keys, buckets, internal bookkeeping and the map's capacity factor. for a large map that becomes substantial. Second, after examining the actual data - I found that 90% of words have no other word in their equivalence group. That means we store a lot of slices containing just a single, usually very short string. it means the slice header itself takes a few times more space then the underlying string it keeps. i.e most of the memory taken by this data structure is wasted on slice headers. In my tests, this data structure used ~60MB of memory to store the relation. Lastly, since each string header and each slice contained in the map adds another pointer hop on the way to the actual data, it could potentially increase cache misses.

MemBlockWordsCatalog - this is a single continuous block of memory(byte slice) containing all strings in the dictionary, stored using one byte per rune(all strings are ascii, so that's possible), where they are sorted first by length and second by lexicographical order of their canonical representations. That means that all words in the same equivalence group are stored right next to each other in the block. We also keep a map indexing from each integer x to the sub slice of strings with length x. When we need to look for a string in the data structure we first look for the sub-slice indexed by the number of runes in that string. if such a sub-slice exists, then we perform binary search inside the sub-slice to locate the start of the string's equivalence group. if such a group is found then we return it, after removing the searched word from the result.

	pros: There's almost no memory waste. The only extra bookkeeping cost is a map[int][]byte holding 31 keys, which is negligible. In my tests the entire data structure needed just over 3MB of memory to store the relation. This is almost a 20X(!!) improvement in memory size compared with the map solution. Also, since the two WordsCatalog implementations are very frequently queried, but never change(after creation), they become good candidates for caching in the cpu. Since the single block implementation takes much less memory there is a good chance large parts of it can fit inside the memory caches of the cpu, leading to much less cache misses. Also, there's much less indirection. all byte slice pointers reference the same continuous block of memory. which means that for the cpu the memory access patterns are easier to handle compared with pointers to many unrelated memory locations.

	cons: the implementation is slightly more complicated. Also binary search typically takes slightly more ops than a map lookup. Lastly, since unlike C, in go strings have headers we have to create a string from the byte slice each time we want to fetch one of the strings out of the block, which takes a tiny bit of extra time for each.


All in all, in my tests the application typically consumed 50MB-60MB less memory when using the MemBlockWordsCatalog implementation (in relative terms that means the application was using 3X-7X more memory when running with the map implementation). Run time was also typically 5%-15% better (faster) with the MemBlockWordsCatalog implementation. Percentage-wise, the difference between the two algorithms tended to increase in favor of the MemBlockWordsCatalog as the test size/amount of requests increased. This was true both for memory consumption and run time. That might be due to the favorable cache warmup associated with the smaller memory size of the MemBlockWordsCatalog implementation.


How to run:
non containerized:
go run "${path_to_app_folder}/main.go" -WordsCatalogAlg="${your_chosen_algorithm_name}"
containerized:
docker build --no-cache -t similarwordsservice:v1 . && docker run --restart=always -p 8000:8000 similarwordsservice:v1 --WordsCatalogAlg="${your_chosen_algorithm_name}"

if you don't specify the WordsCatalogAlg flag you'll get the default implementation, which is the MemBlockWordsCatalog.
you can also just use "docker compose up" but then you cannot define the algorithm to use, and again you get the default one.


Tests, profiling and performance notes:
	There are tests for (almost) every part of the application. The most important of which are the end to end tests involving the test client I built. These tests are found in the test_client package and need a running server to function. you can run them in conjunction with the server profiling test, which runs a server and profiles it, and get extensive profiling information on the server side at the same time.
The results are saved in the 'server_test_results' directory that will be found in the project home dir.
p.s: running the test client against the regular server won't work as it is probing an endpoint which signals server readiness, and the regular server exposes no such endpoint.
to run the server side of the test:
	go test ./server/ -count=1 -v -run TestServerWithProfiling -WordsCatalogAlg="${algorithm}"

to run the client side of the test:
	go test ./test_client -count=1 -v -run "${full_client_test_name}"
	where the client test names can be either TestResponseValidatingTestClient or TestHighLoadTestClient.

** However, the recommended way to run the test client is not directly. **
I've made a bash wrapper script to run it called "run_endtoend_server_test.sh", located in the project's home dir.
the script accepts a single parameter: test name (which must be one of: HighLoadTest ResponseValidatingTest)
it runs the requested test with both words catalog implementations, reports a summary of the results and also displays a comparative plot of the two runs for easier visual inspection.
examle usage:
./run_endtoend_server_test.sh HighLoadTest
HighLoadTest - this test sends in parallel 1000000 similar words requests and 50000 stats requests, and runs a partial validation on the server responses. It is meant to test how the server performs under high stress. When running on my laptop this load is typically handled in ~50 seconds.

ResponseValidatingTest - this test sends and thoroughly validates the correctness of the response for 100000 similar words requests made in parallel on a mix of random strings and words from the dictionary. it is meant to thoroughly check the correctness of the implementation. When running on my laptop this load is typically handled in ~4 seconds.


A few final notes:
- if the server receives a similar words request without a word parameter the response will be: {"error":"word parameter is missing"}

- The search is case insensitive. if you want to change this, just remove the respective part in line 27 of similarwordscontroller.go: 
equivalentWords := wordsCatalog.FindEquivalentWords(strings.ToLower(word))

- I am returning the equivalent words to the search word, even if the searched word was not in the catalog.
For example, the word 'aboil' is in the catalog, but the word 'oabil' is not.
The response for the input 'aboil' will be: { similar:['abilo', 'bailo'] }, while the response for the input 'oabil' will be: { similar:['abilo', 'aboil', 'bailo'] }
The two words contain the same letters, but they receive different responses since the search word needs to be excluded from the result.
