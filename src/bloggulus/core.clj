(ns bloggulus.core)

(defrecord Account [account-id username password email verified])
(defrecord AccountBlog [account-id blog-id])
(defrecord Blog [blog-id feed-url site-url title])
(defrecord Post [post-id blog-id url title preview updated])
(defrecord Session [session-id account-id expiry])

(def db
  {:blogs [(map->Blog {:blog-id 1
                       :feed-url "https://nullprogram.com/feed/"
                       :site-url "https://nullprogram.com/"
                       :title "null program"})
           (map->Blog {:blog-id 2
                       :feed-url "https://eli.thegreenplace.net/feeds/all.atom.xml"
                       :site-url "https://eli.thegreenplace.net/"
                       :title "Eli Bendersky's website"})]
   :posts [(map->Post {:post-id 1
                       :blog-id 1
                       :url "https://nullprogram.com/blog/2021/08/21/"
                       :title "Test cross-architecture without leaving home"
                       :preview "I like to test my software across"
                       :updated #inst "2021-08-21T23:59:33.000-00:00"})
           (map->Post {:post-id 2
                       :blog-id 1
                       :url "https://nullprogram.com/blog/2021/07/30/"
                       :title "strcpy: a niche function you don't need"
                       :preview "The C strcpy function is a common sight"
                       :updated #inst "2021-07-30T19:37:48.000-00:00"})
           (map->Post {:post-id 3
                       :blog-id 2
                       :url "https://eli.thegreenplace.net/2021/plugins-in-go/"
                       :title "Plugins in Go"
                       :preview "Several years ago I started writing a series"
                       :updated #inst "2021-08-28T14:19:00.000-00:00"})
           (map->Post {:post-id 4
                       :blog-id 2
                       :url "https://eli.thegreenplace.net/2021/accessing-postgresql-databases-in-go/"
                       :title "Accessing PostgreSQL databases in Go"
                       :preview "This post discusses some options for accessing"
                       :updated #inst "2021-07-17T13:20:00.000-00:00"})]})

(comment

  (->Blog 1 "http://foo.com/rss.xml" "http://foo.com" "Foo Blog")
  (map->Blog {:blog-id 1
              :feed-url "http://foo.com/rss.xml"
              :site-url "http://foo.com"
              :title "Foo Blog"})

  .)
