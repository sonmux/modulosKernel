#include <linux/module.h>    // kernel modules
#include <linux/kernel.h>    // KERN_INFO
#include <linux/init.h>      // 
#include <linux/seq_file.h>  // var seq_file
#include <linux/proc_fs.h>   // procs file
#include <linux/sched/signal.h> //for process
#include <linux/sched.h> //for process
#include <linux/mm.h> //total ram

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Grupo 5");
MODULE_DESCRIPTION("Cpu");

static void iterate_tasks(struct seq_file *archivo);

struct task_struct *task_padre;       
struct task_struct *task_hijo;       
struct list_head *list;            

struct sysinfo infsys;

//Escribiendo la info
static int write_file(struct seq_file *archivo, void *v)
{   
    seq_printf(archivo, "{\n"); //INI.JSON
    iterate_tasks(archivo);
    seq_printf(archivo, "}\n"); //FIN.JSON
    return 0;
}

//ciclo para obtener los procesos
void iterate_tasks(struct seq_file *archivo) 
{
    //Total de memoria
    long memorytotal;
    si_meminfo(&infsys);
    memorytotal = infsys.totalram * infsys.mem_unit; //bytes
    //Memoria proceso
    unsigned long rss;

    //iteracion formato json
    seq_printf(archivo, " \t\"PROCESOSPADRE\":[\n"); //INI.PROCESOS PADRE

    for_each_process(task_padre){           
        
        seq_printf(archivo, "\t\t{\n"); //INI.PROCESO
        seq_printf(archivo, "\t\t\"PROCESO\":\"%s\",\n", task_padre->comm);
        seq_printf(archivo, "\t\t\"PID\":%d,\n", task_padre->pid);
        seq_printf(archivo, "\t\t\"ESTADO\":%ld,\n", task_padre->state);
        //seq_printf(archivo, "\t\t\"MEMORIA_USO\":%ld\n", task_padre->comm);

        if (task_padre->mm) {
            rss = get_mm_rss(task_padre->mm) << PAGE_SHIFT;
            seq_printf(archivo, "\t\t\"MEMORIA_USO\":%ld,\n", rss);
        } else {
            seq_printf(archivo, "\t\t\"MEMORIA_USO\":0,\n");
        }

        seq_printf(archivo, "\t\t\"MEMORIA_TOTAL\":%ld,\n", memorytotal);
        seq_printf(archivo, "\t\t\"ID_USUARIO\":%d\n", __kuid_val(task_padre->real_cred->uid));
        seq_printf(archivo, "\t\t}"); //FIN.PROCESO

        if(next_task(task_padre) != &init_task){//FIN LISTA
            seq_printf(archivo, ",\n"); //SIG.PROCESOS
        }else{
            seq_printf(archivo, "\n"); //FIN.PROCESOS
        }
    }

    seq_printf(archivo, "\t],\n"); //FIN.PROCESOS PADRE

    //INICIA UN NUEVO OBJETO 

    seq_printf(archivo, " \t\"PROCESOSHIJO\":[\n"); //INI.PROCESOS HIJO

    for_each_process(task_padre){           
        
        list_for_each(list, (&task_padre->children)){ //Iterar procesos padre

            task_hijo = list_entry(list, struct task_struct, sibling); //Obtiene la info del hijo

            seq_printf(archivo, "\t\t{\n"); //INI.PROCESO
            seq_printf(archivo, "\t\t\"PROCESO\":\"%s\",\n", task_hijo->comm);
            seq_printf(archivo, "\t\t\"PID\":%d,\n", task_hijo->pid);
            seq_printf(archivo, "\t\t\"ESTADO\":%ld,\n", task_hijo->state);
            seq_printf(archivo, "\t\t\"PID_PADRE\":%d\n", task_padre->pid);
            seq_printf(archivo, "\t\t},\n"); //FIN.PROCESO
            
          
        }

        if(next_task(task_padre) != &init_task){//FIN LISTA
          
        }else{
            seq_printf(archivo, "\t\t{\n"); //INI.PROCESO
            seq_printf(archivo, "\t\t\"PROCESO\":\"vacio\",\n");
            seq_printf(archivo, "\t\t\"PID\":0,\n");
            seq_printf(archivo, "\t\t\"ESTADO\":0,\n");
            seq_printf(archivo, "\t\t\"PID_PADRE\":0\n");
            seq_printf(archivo, "\t\t}\n"); //FIN.PROCESO
        }

    }

    seq_printf(archivo, "\t]\n"); //FIN.PROCESOS HIHJO
}  

//Se realiza la escritura del archivo
static int to_open(struct inode *inode, struct file *file)
{
    return single_open(file, write_file, NULL);
}

static struct proc_ops operations =
{
    // proc_open --: por error de distr.
	// proc_read --: por error de distr.
	.proc_open = to_open,
    .proc_read = seq_read
    //.open = to_open,
    //.read = seq_read
};


static int iniciar_init(void)
{
    proc_create("cpu_grupo5", 0, NULL, &operations);
    printk(KERN_INFO "M??dulo lista de procesos del Grupo5 Cargado\n");
    return 0;    // 0 = ERROR DE CARGA
}

static void finalizar_end(void)
{
    remove_proc_entry("cpu_grupo5", NULL);
    printk(KERN_INFO "M??dulo lista de procesos del Grupo5 Desmontado\n");
}

module_init(iniciar_init);
module_exit(finalizar_end);