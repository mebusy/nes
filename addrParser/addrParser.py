import sys
import os
import md5 

import sqlite3

if __name__ == '__main__':
    if len(sys.argv) < 2 :
        print 'addrParser <rom_path>'
    rom_path = sys.argv[1]
    
    with open(rom_path) as fp:
        rom = fp.read()
    
    m = md5.new()
    m.update(rom)
    h = m.hexdigest()

    from os.path import expanduser
    home = expanduser("~")

    db_path = os.path.join( home, ".nes/db" , h+".db" )

    print db_path
    conn = sqlite3.connect(db_path)
    c = conn.cursor()

    c.execute("SELECT * FROM address")
    rows = c.fetchall()
    conn.close()

    rows.sort( key = lambda x : int( x[0] , 16 )  )


    mem_ranage = []

    start = None
    end = None
    for addr , _ in rows:
        nAddr  = int( addr  , 16 ) 
        if start is None:
            start = nAddr 
            end = start 
        elif nAddr == end + 1 :
            end = nAddr 
        elif nAddr > end + 1:
            mem_ranage.append(  ( hex( start) , hex(end) , end-start + 1  ) )         
            start = None
            end = None
        else:
            assert False , "data should be error "

    
    mem_ranage.append(  ( hex( start) , hex(end) , end-start + 1  ) )         

    for item in mem_ranage:
        print item 
    print 'done'
