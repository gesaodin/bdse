var opciones = {
        "paging":   true,
        "ordering": true,
        "info":     true,
        "searching": false
    };
//$('#reporte').DataTable(opciones);
$('#reporteSaldos').DataTable(opciones);

$(function(){
    /** 
     * Declaración de variables globales
     */ 
    var rS = $('#reporteSaldos').DataTable();

    $("#cargando").hide();        
    socket.addEventListener('message', function (e) {
        var js = JSON.parse(e.data);
        if (js.tipo != null) {
           switch (js.tipo) {
               case 0:                
                break;
               case 1:
                CrearNotificacion(js.tiempo, js.msj)
                break;
              case 2:
                break;
              case 3:
               ActivarChat(js.De, js.msj, js.tiempo);  
               default:
                   break;
           }
           //console.log(e.data);
        }
        
    });
    socket.addEventListener('open', function (e) {
        console.log("Se establecio la conexión con el socket...");
    });
    socket.addEventListener('close', function (e) {    
        console.log("Se ha cerrado la conexión");    
    });
    
    
    CargarCalendario();    
});

//Listar pagos pendientes del registro de pago.
function LPago(){
    $("#divReporte").show();
    var rS = $('#lstReporte').DataTable();
    var Pago = JSON.stringify({
        agencia: $("#agencia").html()
    });
    rS.clear().draw();
    $.post("api/balance/listarpagos", Pago)
    .done(function (data){    
        var i = 1;
        $.each(data, function(c, v) {
            estatus = v.estatus == null?0:v.estatus;
            switch (estatus) {
                    case 0:
                        estatus = '<span class="label label-warning">Pendiente</span>';
                        break;
                    case 1:
                        estatus = '<span class="label label-success">Procesado</span>';
                        break;
                    case 2:
                        estatus = '<span class="label label-danger">Rechazado</span>';
                        break;    
                    default:
                        //estatus = '<span class="label label-info">Pendiente</span>';
                        break;
                }

            rS.row.add([
                i++,
                v.deposito,                
                v.voucher,
                v.monto,
                estatus
            ]).draw();
        });
    });    


}
function RPago(){
    dep = $("#fechade").val();
    mon = $("#mont").val();
    vouc = $("#vouc").val();
    if(dep == ""){
        $.notify("Debe introducir una fecha ", "warn");
        return false;
    }
    if(mon == ""){
        $.notify("Debe introducir un monto ", "warn");
        return false;
    }

    var Pago = JSON.stringify({
        agencia: $("#agencia").html(),
        voucher: vouc,
        desde : dep,
        hasta : dep,
        deposito: dep,
        banco : parseInt($("#banc").val()),
        formadepago: parseInt($("#form").val()),
        monto: parseFloat(mon),
        observacion: $("#obse").val()
    });
    $.post("api/balance/registrarpago", Pago)
    .done(function (data){    
        $.notify("Proceso exitoso ", "success");
        LFPago();
    });    
}


function RP(){
    $("#mdlRP").modal('show');
}
/**
 * Limpiar el formulario de pago
 */
function LFPago(){
    $("#fechade").val("");
    $("#vouc").val("");
    $("#obse").val("");
    $("#mont").val("");
}
function CargarCalendario(){var a={format:"YYYY/MM/DD",applyLabel:"Aceptar",cancelLabel:"Cancelar",customRangeLabel:"Por Rango",daysOfWeek:["Do","Lu","Ma","Mi","Ju","Vi","Sa"],monthNames:["Enero","Febrero","Marzo","Abril","Mayo","Junio","Julio","Agosto","Septiembre","Octubre","Noviembre","Diciembre"]};$("#daterange-btn").daterangepicker({locale:a,ranges:{Hoy:[moment(),moment()],Ayer:[moment().subtract(1,"days"),moment().subtract(1,"days")],"Hace 7 Dias":[moment().subtract(6,"days"),moment()],"Hace 30 Dias":[moment().subtract(29,"days"),moment()],"Este Mes":[moment().startOf("month"),moment().endOf("month")],"Mes Pasado":[moment().subtract(1,"month").startOf("month"),moment().subtract(1,"month").endOf("month")]},startDate:moment().subtract(29,"days"),endDate:moment()},function(a,b){$("#daterange-btn span").html(a.format("YYYY/MM/DD")+" - "+b.format("YYYY/MM/DD"))}),$("#fechara").daterangepicker({locale:a}),$("#fecha").datepicker({autoclose:!0,format:"yyyy-mm-dd"}),$("#fechade").datepicker({autoclose:!0,format:"yyyy-mm-dd"}),$("#fechadere").datepicker({autoclose:!0,format:"yyyy-mm-dd"})}

function LstSaldoGPS(){
    
    var f = $('#fecharango option:selected').val();
    var suma = 0;

    var f_a = f.split("-");
    if (f_a.length < 3){
        desdeA = f_a[0].replace(" ", "");
        hastaA = f_a[1].replace(" ", "");
        if (f_a[0] == "0"){
            desdeA = moment().format('YYYY/MM/') + '01';
            hastaA = moment().format('YYYY/MM/DD');
        }

        var rfecha =JSON.stringify({
            agencia : $("#agencia").html(),
            desde: desdeA.replace(/\//g, "-"), 
            hasta: hastaA.replace(/\//g, "-")          
        });
        url = "api/balance/cobrosypagos";       
         
        $.post(url,rfecha)
        .done(function(data){ 
           
            PTotales(desdeA, hastaA, data);

        });

       

    }else{
        $.notify("Debe seleccionar un rango", "error");
    }
}

function psFila(fila, buscar){
    var pos = 0;
   
    $.each(fila, function(c, v){
        if (buscar == v){
            pos = c;
        }
    });
    return pos;
}

/**
 * Listar Totales de los saldos para estado de Cuenta
 * 
 */

function PTotales(desde, hasta, data){
    rS = CBTotales();
    var fila = RecorreFechas(desde, hasta, rS);
    var acumuladorSaldosFecha = {};
    var i = 0;
    var total = 0;
     if (data == null){
        
        $.notify("No se encontrarón registros", "error")
        return
    }           
    $.each(data, function(c, v){        
        
        if (i == 0)sAnt = v.vienen == null?0:v.vienen;
        
        i++;
        ingreso = v.ingreso == null?0:v.ingreso;
        egreso = v.egreso == null?0:v.egreso;
        prestamo = v.prestamo == null?0:v.prestamo;
        entregado = v.entregado == null?0:v.entregado;
        recibido = v.recibido == null?0:v.recibido;
        cuota = v.cuota == null?0:v.cuota;
        movimiento = (parseFloat(egreso) + parseFloat(cuota)) - (parseFloat(ingreso) + parseFloat(prestamo)) ;
        saldo = v.saldo == null?0:v.saldo;
        x = parseFloat(entregado) - parseFloat(recibido);
        total = parseFloat(v.saldo) + movimiento + x + sAnt;
        
        fil = psFila(fila,v.fecha);
        rS.cell(fil,1).data(sAnt.toFixed(2)).draw();

        rS.cell(fil,2).data(saldo.toFixed(2)).draw();
        rS.cell(fil,3).data(movimiento.toFixed(2)).draw();
        rS.cell(fil,4).data(entregado.toFixed(2)).draw();
        rS.cell(fil,5).data(recibido.toFixed(2)).draw();

        rS.cell(fil,6).data(total.toFixed(2)).draw();
        sAnt = total;
        
    } );
   
    var table = $('#reporteSaldosGeneral').DataTable();    
        $('#reporteSaldosGeneral tbody').on( 'click', 'tr', function () {
            var data = table.row( this ).data();            
            VentanaEmergente(data[0]);
            
        } );

} 

function CBTotales(){
    $("#divTabla").html('\
    <table class="table table-bordered" id="reporteSaldosGeneral" width="100%">\
        <thead>\
        <tr>\
            <th>Fecha</th>\
            <th>S.Ant</th>\
            <th>S.Día</th>\
            <th>Moviemiento</th>\
            <th>Recibido</th>\
            <th>Entregado</th>\
            <th>Total</th>\
        </tr>\
        </thead>\
    </table>');
    var rS = $('#reporteSaldosGeneral').DataTable(
        {order: [[ 0, "asc" ]]} 
    );

    rS.clear().draw();
    return rS;
}


/**
 * 
 */
function VentanaEmergente(id){
    $('#ventanaEmergenteTitulo').html('Detalle de Sistemas...');
    var tabla = '';
    var data = JSON.stringify({id:0}); 
    var loteria;
    var parley;       
    $.post("api/listasistema",data)
    .done(function(data){      
        tabla = '<table class="table table-bordered" id="reporteSaldosSistemas" width="100%"><thead><tr>';       
        $.each(data, function(c, v){
               tabla += '<th>' + v.nombre + '</th>';
        } );
        tabla += '</tr></thead></table><br><br>';      
        
        $('#ventanaEmergenteContenido').html(tabla);          
        var rS = $('#reporteSaldosSistemas').DataTable( {
            paging:   false,
            ordering: false,
            info:     false,
            searching: false
        });
        //rS.clear().draw();
    });

    
/**
    var data = JSON.stringify({id:1});       
    $.post("api/listasistema",data)
    .done(function(data){      
        tabla = '<table class="table table-bordered" id="reporteSaldosSistemasParley" width="100%"><thead><tr>';       
        $.each(data, function(c, v){
               tabla += '<th>' + v.nombre + '</th>';
        } );
        tabla += '</tr></thead></table>';       
        $('#ventanaEmergenteContenido').append(tabla);        
        var rSs = $('#reporteSaldosSistemasParley').DataTable({
            paging:   false,
            ordering: false,
            info:     false,
            searching: false
        });
        //rSs.clear().draw();
        
    });*/
    console.log( $('#ventanaEmergenteContenido').html() );
    $('#ventanaEmergente').modal({ keyboard: false });   // initialized with no keyboard
    $('#ventanaEmergente').modal('show');                // initializes and invokes show
    $('#ventanaEmergenteContenido').append("HOLA MUNDO");
    
    //PTotalesSistemas(id, $("#agencia").html(), loteria, parley);

   
}

function PTotalesSistemas(fech, agencia, loteria, parley){
    var data = JSON.stringify({
        agencia: agencia, 
        desde:fecha, 
        hasta:fecha
    });
    console.log(loteria);
    loteria.row.add([0,0,0,0,0,0,0,0,0]).draw();
    parley.row.add(['10','HOLA','MUNDO']).draw();
    url = "api/balance/cobrosypagossistemas";   
    $.post(url,data)
    .done(function(data){ 
        //console.log(data);
        $.each(data, function(c, v){
           console.log(v);

        });
        
    });
}

/**
 * @param Date | UNIX
 * @param Date | UNIX
 * @param DataTable
 */
function RecorreFechas(desde, hasta, rS){
    fauxd = desde.split("/");
    fauxh = hasta.split("/");    
    danio = parseInt(fauxd[0]); 
    dmes =  parseInt(fauxd[1]);
    ddia =  parseInt(fauxd[2]);

    hanio = parseInt(fauxh[0]); 
    hmes =  parseInt(fauxh[1]);
    hdia =  parseInt(fauxh[2]);
    var fila = {};
    var count = 0;
    for (h = danio; h <= hanio; h++){
        for (i = dmes; i <= hmes; i++){
            mdmes = new Date(h, i, 0).getDate();
            for (j = ddia; j <= mdmes; j++){
                dia = j;
                if((String(j)).length==1)dia='0'+j;
                mes = i;
                if((String(i)).length==1)mes='0'+i;
                fecha = h + "-" + mes + "-" + dia;
                fila[count] = fecha; 
                 rS.row.add( [
                    fecha,
                    0,
                    0,
                    0,
                    0,
                    0,
                    0,                                               
                ] ).draw( false );
                count++;
                //console.log(danio + "-" + i + "-" + j);
                if(hanio == h && hmes == i && hdia == j )break;                
            }
            ddia = 1;
        } 
    }
    return fila;
}

function btnAccion(id){
    s = '<div class="btn-group">\
        <button type="button" class="btn btn-success">\
        Action</button>\
        <button type="button" class="btn btn-success dropdown-toggle" \
        data-toggle="dropdown" aria-expanded="false">\
        <span class="caret"></span>\
        <span class="sr-only">Toggle Dropdown</span>\
        </button>\
        <ul class="dropdown-menu" role="menu">\
        <li><a href="#">Action</a></li>\
        <li><a href="#">Another action</a></li>\
        <li><a href="#">Something else here</a></li>\
        <li class="divider"></li>\
        <li><a href="#">Separated link</a></li>\
        </ul>\
    </div>';
    return s
}