var opciones = {
        "paging":   true,
        "ordering": true,
        "info":     true,
        "searching": true,
        "language": {
            "decimal": ",",
            "thousands": "."
        }
    };
//$('#reporte').DataTable(opciones);
$('#reporteSaldos').DataTable(opciones);


$(function(){
    /** 
     * Declaración de variables globales
     */
    //var t = $('#reporte').DataTable();    
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
              case 33:
                //Error en los archivos enviados
                CNErr(js.tiempo, js.msj);
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


function CargarCalendario(){
    //Date range picker

    
    var local = {
                "format": 'YYYY/MM/DD',
                "applyLabel": "Aceptar",
                "cancelLabel": "Cancelar",
                "customRangeLabel": 'Por Rango',
                "daysOfWeek": [
                    "Do",
                    "Lu",
                    "Ma",
                    "Mi",
                    "Ju",
                    "Vi",
                    "Sa"
                ],
                "monthNames": [
                    "Enero",
                    "Febrero",
                    "Marzo",
                    "Abril",
                    "Mayo",
                    "Junio",
                    "Julio",
                    "Agosto",
                    "Septiembre",
                    "Octubre",
                    "Noviembre",
                    "Diciembre"
                ],
        };
    //Date range as a button
    $('#daterange-btn').daterangepicker(
        {
         locale : local,
          ranges: {
            'Hoy': [moment(), moment()],
            'Ayer': [moment().subtract(1, 'days'), moment().subtract(1, 'days')],
            'Hace 7 Dias': [moment().subtract(6, 'days'), moment()],
            'Hace 30 Dias': [moment().subtract(29, 'days'), moment()],
            'Este Mes': [moment().startOf('month'), moment().endOf('month')],
            'Mes Pasado': [moment().subtract(1, 'month').startOf('month'), moment().subtract(1, 'month').endOf('month')]
          },
          startDate: moment().subtract(29, 'days'),
          endDate: moment()
        },
        function (start, end) {
          $('#daterange-btn span').html(start.format('YYYY/MM/DD') + ' - ' + end.format('YYYY/MM/DD'));
        }
    );
     
    
    $('#fechara').daterangepicker({locale:local});
    
    $('#fecha').datepicker({autoclose: true, format: 'yyyy-mm-dd'});
    $('#fechade').datepicker({autoclose: true, format: 'yyyy-mm-dd' });
    $('#fechadere').datepicker({autoclose: true, format: 'yyyy-mm-dd' });
}
/**
 * Enviando Archivos
 */
function enviarArchivo(){
      if ($("#archivo").val() == "") {
        $.notify("Debe seleccionar un archivo", "error");
        return false;
      }
      $("#cargando").show();
      
      var formData = new FormData(document.forms.namedItem("forma"));
      $.ajax({
        url: "subirl",
        type: "post",
        dataType: "html",
        data: formData,
        timeout: 15000,
        cache: false,
        contentType: false,
        processData: false
      })
      .done(function(res){        
        $("#archivo").val("");
        $.notify("Envio de archivos exitosos...", "success");
        $("#cargando").hide();
        

      }).fail(function(jqXHR, textStatus){
        $("#archivo").val("");
        if(textStatus === 'timeout'){     
            $.notify("Los archivos exceden el limite en tiempo de conexion intente con menos...", "error");                         
        }
        $("#cargando").hide();
      });
        
}

/**
 * Listar Reporte de Archivos
 * @param t DataTable
 */
function LstRA(){
   
    var f = $('#daterange-btn span').html();
    var f_a = f.split("-");
    if (f_a.length < 3){
        var rfecha =JSON.stringify({
            desde:f_a[0].replace(" ", ""), 
            hasta:f_a[1].replace(" ", "")
        });

        var t = $('#reporte').DataTable();
        t.clear().draw();
        $("#cargando").show();    
        $.post("api/reportearchivo",rfecha)
        .done(function(data){        
            // Get the column API object
            t.column( 1 ).visible( false );
            t.column( 2 ).visible( false );

             $.each(data, function(c, v){
                cantidad  = v.cantidad == null?  0: v.cantidad;
                if(v.estatus == "1")estatus = '<span class="label label-success">Procesado</span>';
                
                switch (cantidad) {
                    case 0:
                        estatus = '<span class="label label-danger">Sin ventas</span>';
                        break;
                    case 1:
                        estatus = '<span class="label label-success">Procesado</span>';
                        break;
                    default:
                        //estatus = '<span class="label label-info">Pendiente</span>';
                        break;
                }
                t.row.add( [
                    parseInt(c) + 1,
                    v.oid,
                    v.tabla,
                    v.nombre,
                    cantidad,
                    v.creado,
                    v.procesado,                    
                    estatus,            
                ] ).draw( false );

            } );
            $("#cargando").hide();
        });
        
        var table = $('#reporte').DataTable();    
        $('#reporte tbody').on( 'click', 'tr', function () {
            var data = table.row( this ).data();
            VentanaEmergente(data[3],"", data[1], data[2]);           
           
        } );
    }else{
        $.notify("Debe seleccionar un rango", "error");
    }
}




/**
 * Crear Ventanas emergentes con diferentes contenidos
 * @param html
 * @param html
 */
function VentanaEmergente(titulo, cont, id, tbl){
    var cont = '<table class="table table-bordered" \
     id="reporteDetallado" width="100%">\
              <thead>\
                <tr>\
                  <th style="width: 5px">#</th>\
                  <th>Agencia</th>\
                  <th >Venta</th>\
                  <th >Premio</th>\
                  <th >Comisión</th>\
                  <th >Saldo</th>\
                </tr>\
              </thead>\
            </table>';
    $('#ventanaEmergenteTitulo').html(titulo);
    $('#ventanaEmergenteContenido').html(cont);

    $('#ventanaEmergente').modal({ keyboard: false })   // initialized with no keyboard
    $('#ventanaEmergente').modal('show')                // initializes and invokes show
    Reporte('reporteDetallado', id, tbl);
}

/**
 * Crear Reportes
 * 
 * @param DataTable
 */
function Reporte(t, oid, tbl){

    var table = $('#' + t).DataTable({
        "paging":   true,
        "ordering": false,
        "info":     false,
        "searching": false
    });   
    var clave = JSON.stringify({id:  oid, tabla: tbl });
    table.clear().draw();
    $.post("api/reportesaldo",clave)
     .done(function(data){         
            $.each(data, function(c, v){
            venta = v.venta == null?  0: v.venta;
            premio = v.premio == null?  0: v.premio;
            comision = v.comision == null?  0: v.comision;
            saldo = v.saldo == null?  0: v.saldo;
            table.row.add( [
                parseInt(c) + 1,
                v.agencia,
                venta,
                premio,
                comision,
                saldo
                            
            ] ).draw( false );

        } );
    });
    
}

/**
 * Crear Listado de Notificaciones
 * 
 * @param string
 * @param string
 */
function CrearNotificacion(t, msj){
    $.notify(msj, "info");
    var cant = parseInt($("#hnoti").html()) + 1;
    $("#hnoti").html(cant)
    $("#tnoti").html("Tienes " + cant + " notificaciones")
    $("#cnoti").append('<li> \
                        <a href="#" title="' + msj + '">\
                        <i class="fa fa-check text-aqua"></i>\
                        ' + msj + '\
                        </a>\
                        </li>');
    

}
function CNErr(t, msj){
    $.notify(msj, "error");
    var cant = parseInt($("#hnoti").html()) + 1;
    $("#hnoti").html(cant)
    $("#tnoti").html("Tienes " + cant + " notificaciones")
    $("#cnoti").append('<li> \
                        <a href="#" title="' + msj + '">\
                        <i class="fa  fa-bug text-red"></i>\
                        ' + msj + '\
                        </a>\
                        </li>');
    

}

/**
 * Activar Mensaje de Chat
 * 
 * @param string
 * @param string
 * @param Date
 */
function ActivarChat(de, msj, t){
    $.notify(msj, "info");
}

/**
 * Listar Saldos Generales de Ventas
 * 
 * @param Date
 */
function LstSaldo(){
    
    var f = $('#daterange-btn span').html();
    var suma = 0;
    var aventa = 0;
    var apremio = 0;
    var acomision = 0;

    var f_a = f.split("-");
    
    if (f_a.length < 3){
        if ($("#tipo option:selected").val() == "0"){
            $.notify("Debe seleccionar un esquema ", "error");
            return 
        }
        var desdeA = f_a[0].replace(" ", "");
        var hastaA = f_a[1].replace(" ", "");

        var rS = $('#reporteSaldos').DataTable();
        var rfecha =JSON.stringify({
            desde:desdeA.replace(/\//g, "-"), 
            hasta:hastaA.replace(/\//g, "-"),
            tabla : $("#tipo option:selected").val(),
            sistema: parseInt($("#sistema option:selected").val())
        });
        rS.clear().draw();
        $("#cargando").show();
       

        $.post("api/reportesaldo",rfecha)
        .done(function(data){        
            $.each(data, function(c, v){
                saldo = v.sal == null?  0: v.sal;
                venta = v.ven == null?  0: v.ven;
                premio = v.pre == null?  0: v.pre;
                comision = v.com == null?  0: v.com;

                suma += parseFloat(saldo);
                aventa += parseFloat(venta);
                apremio += parseFloat(premio);
                acomision += parseFloat(comision);
               
                rS.row.add( [
                    parseInt(c) +1,
                    v.age,
                    venta,
                    premio,
                    comision,
                    saldo,
                    v.fec,                            
                ] ).draw( false );

            } );
            
            $("#tfventa").html(aventa.toFixed(2));
            $("#tfpremio").html(apremio.toFixed(2));
            $("#tfcomision").html(acomision.toFixed(2));
            $("#tfsaldo").html(suma.toFixed(2));
            $("#cargando").hide();
        });
    }else{
        $.notify("Debe seleccionar un rango", "error");
    }
}

function LSistema(){
    var tipo = 0;
    switch ($("#tipo option:selected").val()) {
        case "loteria":
            tipo = 0
            break;
        case "parley":
            tipo = 1;
            break;
        case "todos":
            tipo = 2;
            break;
        default:
            $("#sistema").html("<option value='--'>------------</option>");
            return;
    }
    
    var data = JSON.stringify(
            {
                id:tipo
            }
        );
    
    $.post("api/listasistema",data)
    .done(function(data){
        $("#sistema").html("");        
        $.each(data, function(c, v){
               $("#sistema").append("<option value='" + v.oid + "'>\
               " + v.nombre + "</option>");   

        } );
         $("#sistema").append("<option value=99>Todos</option>");
    });
    
}

function LstSaldoGPS(){
    
    var f = $('#daterange-btn span').html();
    var suma = 0;

    var f_a = f.split("-");
    
    if (f_a.length < 3){
        if ($("#tipo option:selected").val() == "--"){
            $.notify("Debe seleccionar un esquema ", "error");
            return 
        }
        
        var rfecha =JSON.stringify({
            id : parseInt($("#tipo option:selected").val()),
            desde:f_a[0].replace(" ", ""), 
            hasta:f_a[1].replace(" ", "")            
        });
        url = "api/reportesaldogeneral";
        if(parseInt($("#tipo option:selected").val()) == 2) url = "api/balancegeneral";
         
        $.post(url,rfecha)
        .done(function(data){
            switch (parseInt($("#tipo option:selected").val())) {
                case 0:
                    PLoteria(data);
                    break;
                case 1:
                    PParley(data);
                    break;
                case 2: 
                    PTotales(f_a[0].replace(" ", ""), f_a[1].replace(" ", ""), data);
                    break;
                default:
                    break;
            }
            
        });
    }else{
        $.notify("Debe seleccionar un rango", "error");
    }
}
function PLoteria(data){
    $("#divTabla").html('\
        <table class="table table-bordered" id="reporteSaldosGeneral" width="100%">\
            <thead>\
            <tr>\
                <th style="width: 5px">#</th>\
                <th>Fecha</th>\
                <th>Morpheus</th>\
                <th>Pos 1</th>\
                <th>Pos 2</th>\
                <th>Pos 3</th>\
                <th>Maticlo</th>\
                <th>Total</th>\
            </tr>\
            </thead>\
        </table>');
    var rS = $('#reporteSaldosGeneral').DataTable();
    rS.clear().draw();
    var i = 0;
    $.each(data, function(c, v){
        var fila = {};
        var suma = 0;
        $.each(v, function(cl, va){
            saldo = va.saldo == null?  0: va.saldo;
            fila[va.sistema] = saldo;
            suma +=saldo;
        });

        morpheus = fila[1] == null? 0: fila[1]
        pos1 = fila[2] == null? 0: fila[2]
        pos2 = fila[3] == null? 0: fila[3]
        pos3 = fila[4] == null? 0: fila[4]
        maticlo = fila[5] == null? 0: fila[5]
        i++;
        rS.row.add( [
            i,
            c,
            morpheus,
            pos1,
            pos2,
            pos3,
            maticlo,
            suma                                               
        ] ).draw( false );

        
    } );
}

function PParley(data){
    $("#divTabla").html('\
        <table class="table table-bordered" id="reporteSaldosGeneral" width="100%">\
            <thead>\
            <tr>\
                <th style="width: 5px">#</th>\
                <th>Fecha</th>\
                <th>Ilbanquero</th>\
                <th>CyberParley</th>\
                <th>Sport17</th>\
                <th>Total</th>\
            </tr>\
            </thead>\
        </table>');
    var rS = $('#reporteSaldosGeneral').DataTable();
    rS.clear().draw();
    var i = 0;
    $.each(data, function(c, v){
        var fila = {};
        var suma = 0;
        $.each(v, function(cl, va){
            saldo = va.saldo == null?  0: va.saldo;
            fila[va.sistema] = saldo;
            suma +=saldo;
        });

        ilbanquero = fila[6] == null? 0: fila[6]
        cyberparley = fila[7] == null? 0: fila[7]
        sport = fila[8] == null? 0: fila[8]
        
        i++;
        rS.row.add( [
            i,
            c,
            ilbanquero,
            cyberparley,
            sport,
            suma                                               
        ] ).draw( false );

        
    } );
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
    $.each(data, function(c, v){
        //var fila = {};
        var suma = 0;
        var pos = [];
        pos[0] = 2; //Loteria
        pos[1] = 3;  //Parley

        $.each(v, function(cl, va){
            
            //console.log(va);
            saldo = va.saldo == null?  0: va.saldo;
        
            debe = va.debe == null?  0: va.debe;
            haber = va.debe == null?  0: va.haber;
            fil = psFila(fila,cl)-1;
            
            //console.log("Fila: " + fil + " Columna: "+ pos[c] + " Fecha: " + cl);
            rS.cell(fil,pos[c]).data(saldo).draw();

            rS.cell(fil,5).data(debe).draw();
            rS.cell(fil,6).data(haber).draw();
            total = parseFloat(rS.cell(fil,2).data());
            total +=  parseFloat(rS.cell(fil,3).data());
            rS.cell(fil,4).data(total).draw();
               
        });
        total = 0; 
        
    } );


} 

function CBTotales(){
    $("#divTabla").html('\
    <table class="table table-bordered" id="reporteSaldosGeneral" width="100%">\
        <thead>\
        <tr>\
            <th style="width: 5px">#</th>\
            <th>Fecha</th>\
            <th>Loteria</th>\
            <th>Parley</th>\
            <th>Total</th>\
            <th>Entregado</th>\
            <th>Recibido</th>\
            <th>Balance</th>\
        </tr>\
        </thead>\
    </table>');
    var rS = $('#reporteSaldosGeneral').DataTable();
    rS.clear().draw();
    //rS.cells(2,3).data("Hola").draw();
    return rS;
}


function GC(tipo){
    $("#divReporte").html(tableGC());
    $("#reporte").DataTable(opciones);
    var t = $("#reporte").DataTable();
    var data = JSON.stringify({fecha:$("#fecha").val()}); 
    
    $.post("api/balance/cobrosypagos", data)
    .done(function (data){
        
        t.clear().draw();
        var i = 1;
        $.each(data, function(c,v){
            sAnt = v.vienen == null?0:v.vienen;
            ingreso = v.ingreso == null?0:v.ingreso;
            egreso = v.egreso == null?0:v.egreso;
            prestamo = v.prestamo == null?0:v.prestamo;
            entregado = v.entregado == null?0:v.entregado;
            recibido = v.recibido == null?0:v.recibido;
            cuota = v.cuota == null?0:v.cuota;
            movimiento = (parseFloat(egreso) + parseFloat(cuota)) - (parseFloat(ingreso) + parseFloat(prestamo)) ;
            
            x = parseFloat(entregado) - parseFloat(recibido);
            //console.log("SALDO: " + v.saldo + " X: " + x + " MOVIMIENTO : " + movimiento);
            total = parseFloat(v.saldo) + movimiento + x;
            accion = btnAccion(v.agencia, total);
            i++
            if (tipo == 0){
                if(total > 0){
                    t.row.add([
                        accion,
                        v.agencia,
                        sAnt.toFixed(2),
                        parseFloat(v.saldo).toFixed(2),
                        movimiento.toFixed(2),
                        x.toFixed(2),
                        total.toFixed(2)
                    ]).draw();
                }
            }else{
                if(total <= 0){
                    t.row.add([
                        accion,
                        v.agencia,
                        sAnt.toFixed(2),
                        parseFloat(v.saldo).toFixed(2),
                        movimiento.toFixed(2),
                        x.toFixed(2),
                        total.toFixed(2)
                    ]).draw();
                }
            }
            

        })

    })
}
function tableGC(){
    s = '<table class="table table-bordered" cellspacing="0" id="reporte" width="100%">\
            <thead>\
              <tr>\
                <th style="width: 60px">#</th>\
                <th>Agencia</th>\
                <th>S. Ant.</th>\
                <th>S. Día</th>\
                <th>Movimiento</th>\
                <th>+ E - R</th>\
                <th>Total</th>\
              </tr>\
            </thead>\
          </table>';
    return s;

}

function btnAccion(valor,monto){
    s = '<div class="btn-group">\
        <button type="button" class="btn btn-success">\
        <span class="fa fa-cogs"></span></button>\
        <button type="button" class="btn btn-success dropdown-toggle" \
        data-toggle="dropdown" aria-expanded="false">\
        <span class="caret"></span>\
        <span class="sr-only">Toggle Dropdown</span>\
        </button>\
        <ul class="dropdown-menu" role="menu">\
            <li><a href="#" onclick="mdlE(\'mdlMovimiento\',\'\' , \'' + valor + '\', \'' + monto + '\')">Registrar Movimiento</a></li>\
            <li class="divider"></li>\
            <li><a href="#" onclick="mdlE(\'mdlER\',\'er\', \'' + valor + '\', \'' + monto + '\')">Registrar +E -R</a></li>\
            <li><a href="#" onclick="mdlE(\'mdlPre\',\'pre\', \'' + valor + '\', \'' + monto + '\')">Registrar Prestamos </a></li>\
            <li><a href="#" onclick="mdlE(\'mdlEC\',\'ec\', \'' + valor + '\', \'' + monto + '\');">Estado de Cuenta </a></li>\
        </ul>\
    </div>';
    return s
}

function mdlE(id, cod, valor,monto){

    $('#cod' + cod).html(valor);
    $('#' + id).modal('show');
    $('#montoer').val(monto);
    var msj = 'Saldo a cero (0)';
    if(monto > 0){
        msj = 'Saldo a cero (0)'; 
    }
    $('#descripcioner').val(msj);
    $('#divTablaec').html('');
}


function DP(){
    if($("#fecha").val() == ""){
        $.notify("Debe seleccionar la fecha", "error");
        return
    }
    $('#mdlDP').modal('show');
    tabla = $("#rptDeposito").DataTable();
    var data = JSON.stringify({fdeposito:$("#fecha").val()}); 
    $.post("api/movimiento/listardeposito", data)
    .done(function (data){
        tabla.clear().draw();
        var i = 1;
        $.each(data, function(c,v){
           accion = btnADep(v.oid);
            tabla.row.add([
                accion,
                v.agencia,
                v.banco,
                v.voucher,
                v.monto
            ]).draw();
        });
    });
}
function btnADep(oid){
    s = '<div class="btn-group">\
        <button type="button" class="btn btn-success">\
        <span class="fa-check-circle"></span></button>\
        <button type="button" class="btn btn-danger">\
        <span class="fa-times-circle"></span></button></div>';
    return s
}

function EC(){
    $('#cagandoec').show();
    //$("#cagandoec").hide();
    var f = $('#fecharangoec option:selected').val();
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
            agencia : $('#codec').html(),
            desde: desdeA.replace(/\//g, "-"), 
            hasta: hastaA.replace(/\//g, "-")          
        });
        /*
        var rfecha =JSON.stringify({
            agencia :  $('#codec').html(),
            desde:f_a[0].replace(" ", ""), 
            hasta:f_a[1].replace(" ", "")            
        });
        */
        url = "api/balance/cobrosypagos";         
        $.post(url,rfecha)
        .done(function(data){  
                      
            PTotalesDetalles(desdeA, hastaA, data);
            
        });
    }else{
        $.notify("Debe seleccionar un rango", "error");
    }
} 


function EstadoCuenta(){
    $("#divTablaec").html('\
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
    var rS = $('#reporteSaldosGeneral').DataTable({
        scrollY: "200px",
        scrollCollapse: true,
        paging: false,
        searching: true,
        order: [[ 0, "desc" ]]
     });
    rS.clear().draw();
    return rS;
}


function TablaEstadoCuenta(fila, buscar){
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

function PTotalesDetalles(desde, hasta, data){
    rS = EstadoCuenta();
    var fila = RecorreFechas(desde, hasta, rS);
    var acumuladorSaldosFecha = {};
    var i = 0;
    var total = 0;
    if (data == null){        
        $.notify("No se encontrarón registros", "error")
        return
    } 
    $.each(data, function(c, v){        
        sAnt = v.vienen == null?0:v.vienen;
        ingreso = v.ingreso == null?0:v.ingreso;
        egreso = v.egreso == null?0:v.egreso;
        prestamo = v.prestamo == null?0:v.prestamo;
        entregado = v.entregado == null?0:v.entregado;
        recibido = v.recibido == null?0:v.recibido;
        cuota = v.cuota == null?0:v.cuota;
        movimiento = (parseFloat(egreso) + parseFloat(cuota)) - (parseFloat(ingreso) + parseFloat(prestamo)) ;
        saldo = v.saldo == null?0:v.saldo;
        x = parseFloat(entregado) - parseFloat(recibido);
        total = parseFloat(v.saldo) + movimiento + x;
        
        fil = psFila(fila,v.fecha);
        rS.cell(fil,1).data(sAnt.toFixed(2)).draw();

        rS.cell(fil,2).data(saldo.toFixed(2)).draw();
        rS.cell(fil,3).data(movimiento.toFixed(2)).draw();
        rS.cell(fil,4).data(entregado.toFixed(2)).draw();
        rS.cell(fil,5).data(recibido.toFixed(2)).draw();

        rS.cell(fil,6).data(total.toFixed(2)).draw();

        
    } );


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


function RegistrarER(){
    var EntregadoRecibido = JSON.stringify ({
        agencia: $("#coder").html(),
        fecha : $("#fechadere").val(),
        deposito : $("#fechadere").val(),
        forma: parseInt($("#tipoer option:selected").val()), //0 Entregado: DEBE 1 Recibido:HABER
        banco: parseInt($("#cuentaer option:selected").val()),
        monto: parseFloat($("#montoer").val()),
        voucher: $("#voucer").val(),
        observacion: $("#descripcioner").val(),
        estatus: 1
    });

    url = "api/balance/registrarpago";         
    $.post(url,EntregadoRecibido)
    .done(function(data){  
        $('#mdlER').modal('hide');
        $.notify("Registro Exitoso...", "success");
    });
    
}