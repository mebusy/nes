from PIL import Image , ImageDraw
import struct

def showRegion( data ) :
    w = 24  # world map width
    h = 20  # ? 
    tilew = 8
    tileh = 8
    im = Image.new( "RGB", ( w*tilew, h*tileh  ) )
    draw = ImageDraw.Draw(im)
    cols = ( "#000000" ,"#46BAFE" ,"#434E29" , "#60F526", "#6D26A4", "#C2F618","#C34C31" , "#382DF7" , "#A623FC"  )
    output = ''
    for y in xrange( h ):
        for x in xrange( w/2 ):
            byte =  data[ y * ( w/2) + x  ] 
            d = struct.unpack( "B" , byte)[0]
            indices = [ d >> 4 , d & 0xF   ]
            for i, idx in enumerate(indices): 
                output += str(idx)
                col = cols[ idx ]
                draw.rectangle( [ (2*x+i)*tilew, (y) * tileh ,  (2*x+i+1)*tilew, (y+1) * tileh ] , fill = col )

        output += '\n'

    print output

    del draw 
    im.show()
    
                
            



if __name__  == "__main__":
    with open( "./doaeU.nes" , "rb" ) as fp:
        data = fp.read()

    addr = 0x37df5
    showRegion( data[addr: addr + 24*20] ) 

    print 'done'
